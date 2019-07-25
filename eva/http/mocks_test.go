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
)

type mockSource struct {
}

// Connections to the eva service.
func (source *mockSource) Tenant() eva.Tenant {
	return nil
}

// Connection returns a specific connection channel based on the category of data.
func (source *mockSource) Connection(label interface{}) (channel eva.ConnectionChannel, err error) {
	return nil, nil
}

// LatestSnapshot returns the latest snapshot channel for the specified category of data.
func (source *mockSource) LatestSnapshot(label interface{}) (channel eva.SnapshotChannel, err error) {
	return nil, nil
}

// AsOfSnapshot returns the snapshot channel of the specified category as of the rules provided.
func (source *mockSource) AsOfSnapshot(label interface{}, asOf interface{}) (channel eva.SnapshotChannel, err error) {
	return nil, nil
}

// Query the source for data.
func (source *mockSource) Query(query interface{}, parameters ...interface{}) (result eva.Result, err error) {
	return nil, nil
}

// CanLog checks if the logger can log.
func (source *mockSource) Serializer() (edn.Serializer, error) {
	return edn.EvaEdnMimeType, nil
}
