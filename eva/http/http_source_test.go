// Copyright 2018-2019 Workiva Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Workiva/eva-client-go/edn"
	"github.com/Workiva/eva-client-go/eva"
	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeClient struct {
	contentType string
	status      int
}

type fakeCaller struct {
	tries     int
	callCount int
	minTime   time.Duration
	t         time.Time
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		Status:     "Testing",
		StatusCode: c.status,
		Header: map[string][]string{
			"test": {"val"},
		},
		Body: ioutil.NopCloser(strings.NewReader("[]")),
	}

	if len(c.contentType) > 0 {
		resp.Header.Add("Content-Type", c.contentType)
	}

	return resp, nil
}

func fakeRetryCaller(tries int) *fakeCaller {
	return &fakeCaller{
		tries:     tries,
		minTime:   time.Millisecond * 100000,
		callCount: 0,
		t:         time.Now().Add(-time.Hour),
	}
}

func (fake *fakeCaller) clientFunc(c httpDoer, r *http.Request) (resp *http.Response, err error) {

	now := time.Now()

	dur := now.Sub(fake.t) / time.Millisecond
	if dur < fake.minTime {
		fake.minTime = dur
	}
	fake.t = now
	fake.callCount++

	if fake.callCount == fake.tries {
		f := &fakeClient{
			status: http.StatusOK,
		}
		resp, err = f.Do(r)
	} else {
		err = &url.Error{
			Err: fmt.Errorf("EOF"),
		}
	}

	return resp, err
}

func fakeGoodCaller(contentType string) func(c httpDoer, r *http.Request) (*http.Response, error) {
	return func(c httpDoer, r *http.Request) (*http.Response, error) {
		f := &fakeClient{
			status:      http.StatusOK,
			contentType: contentType,
		}
		return f.Do(r)
	}
}

var (
	fakeBadCaller = func(c httpDoer, r *http.Request) (*http.Response, error) {
		f := &fakeClient{
			status: http.StatusTeapot,
		}
		return f.Do(r)
	}
)

var _ = Describe("Binding Test", func() {
	Context("with the default marshaller", func() {
		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost",
					"cert":   "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			snap, err := source.AsOfSnapshot("test", 123)
			Ω(err).Should(BeNil())
			Ω(snap).ShouldNot(BeNil())
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"retries": "42",
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSrc, is := source.(*httpSourceImpl); is {
				Ω(httpSrc.retryTimes).Should(BeEquivalentTo(42))
				Ω(httpSrc.retryPause).Should(BeEquivalentTo(time.Duration(defaultRetryPauseTimeout) * time.Millisecond))
			} else {
				Fail("Expecting HTTP Source")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"retries": "@99",
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).ShouldNot(BeNil())
			Ω(source).Should(BeNil())
			Ω(err).Should(test.HaveMessage(eva.ErrInvalidConfiguration))
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"retries": "42@",
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).ShouldNot(BeNil())
			Ω(source).Should(BeNil())
			Ω(err).Should(test.HaveMessage(eva.ErrInvalidConfiguration))
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"retries": "42@99",
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSrc, is := source.(*httpSourceImpl); is {
				Ω(httpSrc.retryTimes).Should(BeEquivalentTo(42))
				Ω(httpSrc.retryPause).Should(BeEquivalentTo(time.Duration(99) * time.Millisecond))
			} else {
				Fail("Expecting HTTP Source")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost",
					"cert":   "badcert"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).ShouldNot(BeNil())
			Ω(source).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidCertificate))
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).ShouldNot(BeNil())
			Ω(source).Should(BeNil())
			Ω(err).Should(test.HaveMessage(eva.ErrInvalidConfiguration))
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient(&fakeClient{
					status: http.StatusOK,
				}, nil)
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			_, is := source.(*httpSourceImpl)
			Ω(is).Should(BeTrue())
		})

		It("compile the wildcard pattern correctly", func() {
			source, err := initHttpSource(nil, nil)
			Ω(err).ShouldNot(BeNil())
			Ω(source).Should(BeNil())
			Ω(err).Should(test.HaveMessage(eva.ErrInvalidConfiguration))
		})

		It("construct the url", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				url := httpSource.formulateUrl("oper")
				Ω(url).Should(BeEquivalentTo("http://localhost/eva/v.1/oper/tenant/test"))
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("construct the url", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewCorrelationTenant("tenant", "foo")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = nil

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(test.HaveMessage(ErrNoServiceImpl))
				Ω(res).Should(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewCorrelationTenant("tenant", "foo")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())

				str, h := res.String()
				Ω(h).Should(BeTrue())
				Ω(str).Should(BeEquivalentTo("[]"))
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewCorrelationTenant("tenant", "foo")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())

				var has bool
				err, has = res.Error()
				Ω(err).Should(BeNil())
				Ω(has).Should(BeFalse())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewCorrelationTenant("tenant", "foo")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())

				// mess with the result.
				res.(*httpResult).examine = nil

				var has bool
				err, has = res.Error()
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(test.HaveMessage(eva.ErrInvalidSerializer))
				Ω(has).Should(BeTrue())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewCorrelationTenant("tenant", "foo")
			Ω(err).Should(BeNil())

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration

			timeout := 100
			tries := 5
			config, err = eva.NewConfiguration(fmt.Sprintf(`{
				"source": {
					"type":   "http",
					"server": "localhost",
					"retries": "%d@%d"
				},
				"category": "test"
			}`, tries, timeout))
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewCorrelationTenant("tenant", "foo")
			Ω(err).Should(BeNil())

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				f := fakeRetryCaller(tries)
				httpSource.callClient = f.clientFunc

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())
				Ω(f.callCount).Should(BeEquivalentTo(tries))
				Ω(f.minTime).Should(BeNumerically(">=", timeout))
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost",
					"cert":   "-----BEGIN CERTIFICATE-----\nMIICEjCCAXsCAg36MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG\nA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE\nMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl\nYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw\nODIyMDUyNjU0WhcNMTcwODIxMDUyNjU0WjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE\nCAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs\nZS5jb20wXDANBgkqhkiG9w0BAQEFAANLADBIAkEAm/xmkHmEQrurE/0re/jeFRLl\n8ZPjBop7uLHhnia7lQG/5zDtZIUC3RVpqDSwBuw/NTweGyuP+o8AG98HxqxTBwID\nAQABMA0GCSqGSIb3DQEBBQUAA4GBABS2TLuBeTPmcaTaUW/LCB2NYOy8GMdzR1mx\n8iBIu2H6/E2tiY3RIevV2OW61qY2/XRQg7YPxx3ffeUugX9F4J/iPnnu1zAxxyBy\n2VguKv4SWjRFoRkIfIlHX0qVviMhSlNy2ioFLy7JcPZb+v3ftDGywUqcBiVDoea0\nHn+GmxZA\n-----END CERTIFICATE-----"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())

			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeBadCaller

				form := url.Values{}

				form.Add("foo", "bar")
				res, err := httpSource.call("GET", "http://localhost", form)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())

				var has bool
				err, has = res.Error()
				Ω(err).Should(test.HaveMessage(ErrServiceError))
				Ω(has).Should(BeTrue())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				foo, err := edn.NewStringElement("foo")
				Ω(err).Should(BeNil())
				res, err := httpSource.queryImpl(foo)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())

				res, err := httpSource.queryImpl("\"foo\"")
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())
				tenant, err := eva.NewTenant("tenant")
				Ω(err).Should(BeNil())

				res, err := httpSource.queryImpl(tenant, 42)
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(test.HaveMessage(ErrUnsupportedType))
				Ω(res).Should(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("form with no params", func() {
			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			err = source.(*httpSourceImpl).fillForm(url.Values{})
			Ω(err).Should(BeNil())
		})

		It("form with string params", func() {
			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			err = source.(*httpSourceImpl).fillForm(url.Values{}, "foo")
			Ω(err).Should(BeNil())
		})

		It("form with bad params", func() {
			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			err = source.(*httpSourceImpl).fillForm(url.Values{}, &struct{}{})
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnsupportedType))
		})

		It("form with elem params", func() {
			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			foo, err := edn.NewStringElement("foo")
			Ω(err).Should(BeNil())

			err = source.(*httpSourceImpl).fillForm(url.Values{}, foo)
			Ω(err).Should(BeNil())
		})

		It("form with ref params", func() {
			ref, err := eva.NewConnectionReference("label")
			Ω(err).Should(BeNil())

			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var srcConfig eva.SourceConfiguration
			srcConfig, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(srcConfig).ShouldNot(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			err = source.(*httpSourceImpl).fillForm(url.Values{}, ref)
			Ω(err).Should(BeNil())
		})
	})
})
