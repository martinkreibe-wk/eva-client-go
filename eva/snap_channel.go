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

import (
	"github.com/Workiva/eva-client-go/edn"
)

const (
	// SnapshotReferenceType defines a new snapshot reference.
	SnapshotReferenceType ChannelType = "eva.client.service/snapshot-ref"
)

// SnapshotChannel defines the channel to a particular eva snapshot
type SnapshotChannel interface {
	Channel

	// Pull from the snapshot.
	Pull(pattern interface{}, ids interface{}, parameters ...interface{}) (Result, error)

	// Invoke from the snapshot
	Invoke(function interface{}, parameters ...interface{}) (Result, error)

	// AsOf the time specified.
	AsOf() *int
}

type PullImplementation func(pattern edn.Serializable, ids edn.Serializable, params ...interface{}) (result Result, err error)
type InvokeImplementation func(function edn.Serializable, parameters ...interface{}) (result Result, err error)

type BaseSnapshotChannel struct {
	*BaseChannel
	pullImpl   PullImplementation
	invokeImpl InvokeImplementation
}

func NewBaseSnapshotChannel(label edn.Serializable, source Source, pullImpl PullImplementation, invokeImpl InvokeImplementation, asOf interface{}) (channel *BaseSnapshotChannel, err error) {

	var asOfSer edn.Serializable
	if asOfSer, err = decodeSerializable(asOf); err == nil {
		var base *BaseChannel
		if base, err = NewBaseChannel(
			SnapshotReferenceType,
			source, map[string]edn.Serializable{
				LabelReferenceProperty: label,
				AsOfReferenceProperty:  asOfSer,
			}); err == nil {
			channel = &BaseSnapshotChannel{
				BaseChannel: base,
				pullImpl:    pullImpl,
				invokeImpl:  invokeImpl,
			}
		}
	}

	return channel, err
}

// Label to this particular channel
func (channel *BaseSnapshotChannel) Label() string {
	return channel.BaseChannel.Label()
}

// AsOf the time specified.
func (channel *BaseSnapshotChannel) AsOf() (asOf *int) {

	if asOfSer := channel.Reference().GetProperty(AsOfReferenceProperty); asOfSer != nil {
		switch val := asOfSer.(type) {
		case edn.Element:
			if val.ElementType() == edn.IntegerType {
				v := int(val.Value().(int64))
				asOf = &v
			}
		case rawIntImpl:
			v := int(val.Int())
			asOf = &v
		}
	}

	return asOf
}

// Pull from the snapshot.
func (channel *BaseSnapshotChannel) Pull(pattern interface{}, ids interface{}, parameters ...interface{}) (result Result, err error) {

	var ptrn edn.Serializable
	var idSer edn.Serializable

	if ptrn, err = decodeSerializable(pattern); err == nil {
		idSer, err = decodeSerializable(ids)
	}

	if err == nil {
		result, err = channel.pullImpl(ptrn, idSer, parameters...)
	}

	return result, err
}

// Invoke from the snapshot
func (channel *BaseSnapshotChannel) Invoke(function interface{}, parameters ...interface{}) (result Result, err error) {

	var funcElem edn.Serializable
	funcElem, err = decodeSerializable(function)

	if err == nil {
		result, err = channel.invokeImpl(funcElem, parameters...)
	}

	return result, err
}
