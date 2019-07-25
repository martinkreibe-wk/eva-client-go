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

func newHttpSnapChannel(connChan *httpConnChanImpl, t edn.Element) (channel eva.SnapshotChannel, err error) {

	snap := &httpSnapChanImpl{
		connChan: connChan,
	}

	var base *eva.BaseSnapshotChannel
	label, err := connChan.Reference().GetProperty(eva.LabelReferenceProperty)
	if err != nil {
		return nil, err
	}

	if base, err = eva.NewBaseSnapshotChannel(label, connChan.Source(), snap.pull, snap.invoke, t); err == nil {
		snap.BaseSnapshotChannel = base
		channel = snap
	}

	return channel, err
}

func (snap *httpSnapChanImpl) invoke(function edn.Element, parameters ...interface{}) (result eva.Result, err error) {
	uri := snap.connChan.Source().(*httpSourceImpl).formulateUrl("invoke")

	var serializer edn.Serializer
	if serializer, err = snap.connChan.Source().(*httpSourceImpl).Serializer(); err == nil {
		form := url.Values{}
		if err = snap.connChan.Source().(*httpSourceImpl).fillForm(form, parameters...); err == nil {

			stream := edn.NewStringStream()
			err = serializer.SerializeTo(stream, function)
			if err != nil {
				return nil, err
			}

			form.Set("function", stream.String())

			if ref := snap.Reference(); ref != nil {

				stream = edn.NewStringStream()
				err = serializer.SerializeTo(stream, ref)
				if err != nil {
					return nil, err
				}

				form.Add("reference", stream.String())
				result, err = snap.connChan.Source().(*httpSourceImpl).call(http.MethodPost, uri, form)
			}
		}
	}

	return result, err
}

func (snap *httpSnapChanImpl) pull(pattern edn.Element, ids edn.Element, params ...interface{}) (result eva.Result, err error) {

	uri := snap.connChan.Source().(*httpSourceImpl).formulateUrl("pull")
	form := url.Values{}

	var serializer edn.Serializer
	if serializer, err = snap.connChan.Source().(*httpSourceImpl).Serializer(); err == nil {

		stream := edn.NewStringStream()
		err = serializer.SerializeTo(stream, pattern)
		if err != nil {
			return nil, err
		}

		form.Add("pattern", stream.String())

		stream = edn.NewStringStream()
		err = serializer.SerializeTo(stream, ids)
		if err != nil {
			return nil, err
		}

		form.Add("ids", stream.String())

		stream = edn.NewStringStream()
		err = serializer.SerializeTo(stream, snap.Reference())
		if err != nil {
			return nil, err
		}

		form.Add("reference", stream.String())
	}

	if err == nil {
		err = snap.connChan.Source().(*httpSourceImpl).fillForm(form, params...)
	}

	if err == nil {
		result, err = snap.connChan.Source().(*httpSourceImpl).call(http.MethodPost, uri, form)
	}

	return result, err
}
