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
	"github.com/Workiva/eva-client-go/edn"
	"github.com/Workiva/eva-client-go/eva"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (

	// ErrServiceError defines a service error.
	ErrServiceError = edn.ErrorMessage("Service call failed")
)

type httpResult struct {
	body        []byte
	code        int
	contentType string
	examine     eva.ErrorExaminer
}

func newHttpResult(req *http.Request, form url.Values, resp *http.Response) (result eva.Result, err error) {

	var data []byte
	if resp.Body != nil {
		data, err = ioutil.ReadAll(resp.Body)
		err = edn.AppendError(err, resp.Body.Close())
	}

	log(req, form, resp, data)

	var serializer edn.Serializer
	var contentType string
	if contentType = resp.Header.Get("Content-Type"); len(contentType) > 0 {
		serializer, err = edn.GetSerializer(contentType)
	} else {
		serializer = edn.DefaultMimeType
	}

	var examiner eva.ErrorExaminer
	if err == nil {
		examiner, err = eva.GetErrorExaminer(serializer)
	}

	if err == nil {
		result = &httpResult{
			body:        data,
			code:        resp.StatusCode,
			contentType: contentType,
			examine:     examiner,
		}
	}

	return result, err
}

func log(req *http.Request, form url.Values, resp *http.Response, body []byte) {
	//fmt.Printf("---- New Request ----\n")
	//fmt.Printf("\tMethod: %s\n", req.Method)
	//fmt.Printf("\tURI: %s\n", req.URL)
	//for n, v := range req.Header {
	//	fmt.Printf("\t\tHeader [%s]: %s\n", n, v)
	//}
	//
	//for n, v := range form {
	//	fmt.Printf("\t\tForm [%s]: %s\n", n, v)
	//}
	//
	//fmt.Printf("---- New Response ----\n")
	//fmt.Printf("\tStatus Code: %v\n", resp.StatusCode)
	//for n, v := range resp.Header {
	//	fmt.Printf("\t\tHeader [%s]: %s\n", n, v)
	//}
	//
	//fmt.Printf("\tBody: %s\n", string(body))
}

// String this result.
func (result *httpResult) String() (string, bool) {
	body := string(result.body)
	return body, len(body) > 0
}

// Error of this result.
func (result *httpResult) Error() (err error, _ bool) {

	if result.examine != nil {
		if result.code < http.StatusOK || result.code >= http.StatusBadRequest {
			err = edn.MakeError(ErrServiceError, result)
		}

		if err == nil {
			err = result.examine(result.body)
		}
	} else {
		err = edn.MakeErrorWithFormat(eva.ErrInvalidSerializer, "Unsupported return type: %s", result.contentType)
	}

	return err, err != nil
}
