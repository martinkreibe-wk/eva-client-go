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

type ChannelType string

// Channel defines an underlying communication mechanism to an eva construct.
type Channel interface {

	// Label to this particular channel
	Label() (string, error)

	// Type of channel.
	Type() ChannelType

	// Reference EDN of this node.
	Reference() Reference

	// Source this channel connects to
	Source() Source
}
