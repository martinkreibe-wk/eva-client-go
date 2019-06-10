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

const (
	// ConnectionReferenceType defines a new connection reference.
	ConnectionReferenceType ChannelType = "eva.client.service/connection-ref"
)

// ConnectionChannel defines the channel to the eva connection
type ConnectionChannel interface {
	Channel

	// Transact the data to the channel
	Transact(data ...interface{}) (Result, error)

	// LatestSnapshot returns the latest snapshot channel.
	LatestSnapshot() (SnapshotChannel, error)

	// AsOfSnapshot returns the snapshot channel as of the rules provided.
	AsOfSnapshot(interface{}) (SnapshotChannel, error)
}
