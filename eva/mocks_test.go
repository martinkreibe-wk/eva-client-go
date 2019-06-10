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

package eva

import "github.com/Workiva/eva-client-go/edn"

type mockResult struct {
}

// String version of the call.
func (mock *mockResult) String() (string, bool) {
	return "test", true
}

// Error from the call.
func (mock *mockResult) Error() (error, bool) {
	return nil, false
}

type mockSource struct {
}

// Connections to the eva service.
func (source *mockSource) Tenant() Tenant {
	return nil
}

// Connection returns a specific connection channel based on the category of data.
func (source *mockSource) Connection(category interface{}) (channel ConnectionChannel, err error) {
	return nil, nil
}

// LatestSnapshot returns the latest snapshot channel for the specified category of data.
func (source *mockSource) LatestSnapshot(category interface{}) (channel SnapshotChannel, err error) {
	return nil, nil
}

// AsOfSnapshot returns the snapshot channel of the specified category as of the rules provided.
func (source *mockSource) AsOfSnapshot(category interface{}, asOf interface{}) (channel SnapshotChannel, err error) {
	return nil, nil
}

// Query the source for data.
func (source *mockSource) Query(query interface{}, parameters ...interface{}) (result Result, err error) {
	return nil, nil
}

// CanLog checks if the logger can log.
func (source *mockSource) Serializer() (edn.Serializer, error) {
	return nil, nil
}

func mockQuery(_ interface{}, _ ...interface{}) (Result, error) {
	return &mockResult{}, nil
}

func makeMockConnChannel(label edn.Serializable, source Source) (c ConnectionChannel, e error) {
	return NewBaseConnectionChannel(
		label,
		source,
		func(transaction edn.Serializable) (Result, error) {
			return &mockResult{}, nil
		},
		func(asOf edn.Serializable) (SnapshotChannel, error) {
			return NewBaseSnapshotChannel(
				label,
				source,
				func(pattern edn.Serializable, ids edn.Serializable, params ...interface{}) (result Result, err error) {
					return &mockResult{}, nil
				},
				func(function edn.Serializable, parameters ...interface{}) (result Result, err error) {
					return &mockResult{}, nil
				},
				asOf)
		})
}
