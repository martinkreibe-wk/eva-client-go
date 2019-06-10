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

const (

	// ErrInvalidSource defines an error with the source
	ErrInvalidSource = edn.ErrorMessage("Invalid source")
)

// Channel defines an underlying communication mechanism to an eva construct.
type BaseChannel struct {
	ref    Reference
	source Source
}

// NewBaseChannel creates a new base channel from the channel type, the label and the source.
func NewBaseChannel(refType ChannelType, source Source, properties map[string]edn.Serializable) (channel *BaseChannel, err error) {
	if source != nil {
		if ref, err := newReference(refType, properties); err == nil {
			channel = &BaseChannel{
				ref:    ref,
				source: source,
			}
		}
	} else {
		err = edn.MakeError(ErrInvalidSource, nil)
	}

	return channel, err
}

// Label to this particular channel
func (channel *BaseChannel) Label() (label string) {
	ser := channel.ref.GetProperty(LabelReferenceProperty)
	if ser != nil {
		switch val := ser.(type) {
		case edn.Element:
			if val.ElementType() == edn.StringType {
				label = val.Value().(string)
			} else {
				label = val.String()
			}
		case *rawStringImpl, rawStringImpl:
			label = val.String()
		}
	}

	return label
}

// Type of channel.
func (channel *BaseChannel) Type() ChannelType {
	return channel.ref.Type()
}

// Reference of this channel.
func (channel *BaseChannel) Reference() Reference {
	return channel.ref
}

// Source this channel connects to
func (channel *BaseChannel) Source() Source {
	return channel.source
}
