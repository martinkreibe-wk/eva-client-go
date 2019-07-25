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
type httpConnChanImpl struct {
	*eva.BaseConnectionChannel
}

// newHttpConnChannel creates a new connection channel.
func newHttpConnChannel(label edn.Element, source eva.Source) (channel eva.ConnectionChannel, err error) {

	httpConn := &httpConnChanImpl{}

	var base *eva.BaseConnectionChannel
	if base, err = eva.NewBaseConnectionChannel(label, source, httpConn.transact, httpConn.asOfSnapshot); err == nil {
		httpConn.BaseConnectionChannel = base
		channel = httpConn
	}

	return channel, err
}

// asOfSnapshot is the implementation for getting a snapshot at a particular reference.
func (connChan *httpConnChanImpl) asOfSnapshot(t edn.Element) (channel eva.SnapshotChannel, err error) {
	return newHttpSnapChannel(connChan, t)
}

// transact will transact an edn to the eva database.
// Submits a transaction, blocking until a result is available.
func (connChan *httpConnChanImpl) transact(transaction edn.Element) (result eva.Result, err error) {
	form := url.Values{}

	var serializer edn.Serializer
	if serializer, err = connChan.Source().Serializer(); err != nil {
		return nil, err
	}

	if ref := connChan.Reference(); ref != nil {
		stream := edn.NewStringStream()
		err = serializer.SerializeTo(stream, ref)
		if err != nil {
			return nil, err
		}

		form.Add("reference", stream.String())
	}

	stream := edn.NewStringStream()
	err = serializer.SerializeTo(stream, transaction)
	if err != nil {
		return nil, err
	}

	form.Add("transaction", stream.String())

	switch source := connChan.Source().(type) {
	case *httpSourceImpl:
		uri := source.formulateUrl("transact")
		result, err = source.call(http.MethodPost, uri, form)
	default:
		err = edn.MakeErrorWithFormat(ErrUnsupportedType, "source type: %T", source)
	}

	return result, err
}
