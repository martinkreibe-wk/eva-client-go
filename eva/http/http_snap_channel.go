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
	"net/http"
	"net/url"
)

// httpConnChanImpl defines the connection channel for the http source.
type httpSnapChanImpl struct {
	*eva.BaseSnapshotChannel
	connChan *httpConnChanImpl
}

func newHttpSnapChannel(connChan *httpConnChanImpl, t edn.Serializable) (channel eva.SnapshotChannel, err error) {

	snap := &httpSnapChanImpl{
		connChan: connChan,
	}

	var base *eva.BaseSnapshotChannel
	if base, err = eva.NewBaseSnapshotChannel(connChan.Reference().GetProperty(eva.LabelReferenceProperty), connChan.Source(), snap.pull, snap.invoke, t); err == nil {
		snap.BaseSnapshotChannel = base
		channel = snap
	}

	return channel, err
}

func (snap *httpSnapChanImpl) invoke(function edn.Serializable, parameters ...interface{}) (result eva.Result, err error) {
	uri := snap.connChan.Source().(*httpSourceImpl).formulateUrl("invoke")

	var serializer edn.Serializer
	if serializer, err = snap.connChan.Source().(*httpSourceImpl).Serializer(); err == nil {
		form := url.Values{}
		if err = snap.connChan.Source().(*httpSourceImpl).fillForm(form, parameters...); err == nil {
			var query string
			if query, err = function.Serialize(serializer); err == nil {
				form.Set("function", query)

				if ref := snap.Reference(); ref != nil {
					var str string
					if str, err = ref.Serialize(serializer); err == nil {
						form.Add("reference", str)
						result, err = snap.connChan.Source().(*httpSourceImpl).call(http.MethodPost, uri, form)
					}
				}
			}
		}
	}

	return result, err
}

func (snap *httpSnapChanImpl) pull(pattern edn.Serializable, ids edn.Serializable, params ...interface{}) (result eva.Result, err error) {

	uri := snap.connChan.Source().(*httpSourceImpl).formulateUrl("pull")
	form := url.Values{}

	var serializer edn.Serializer
	if serializer, err = snap.connChan.Source().(*httpSourceImpl).Serializer(); err == nil {

		if err == nil {

			var ptrnStr string
			if ptrnStr, err = pattern.Serialize(serializer); err == nil {
				form.Add("pattern", ptrnStr)

				var idsStr string
				if idsStr, err = ids.Serialize(serializer); err == nil {
					form.Add("ids", idsStr)
				}
			}

			if err == nil {
				if ref := snap.Reference(); ref != nil {
					var str string
					if str, err = ref.Serialize(serializer); err == nil {
						form.Add("reference", str)
					}
				}
			}
		}
	}

	if err == nil {
		err = snap.connChan.Source().(*httpSourceImpl).fillForm(form, params...)
	}

	if err == nil {
		result, err = snap.connChan.Source().(*httpSourceImpl).call(http.MethodPost, uri, form)
	}

	return result, err
}
