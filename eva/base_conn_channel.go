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

// ConnectionChannel defines the channel to the eva connection
type BaseConnectionChannel struct {
	*BaseChannel
	transactImpl      TransactImpl
	asOfSnapshotImpl  AsOfSnapshotImpl
	tenantsWithSchema map[string]bool
}

type AsOfSnapshotImpl func(asOf edn.Serializable) (SnapshotChannel, error)
type TransactImpl func(transaction edn.Serializable) (Result, error)

func NewBaseConnectionChannel(label edn.Serializable, source Source, transactImpl TransactImpl, asOfSnapshotImpl AsOfSnapshotImpl) (channel *BaseConnectionChannel, err error) {

	if label != nil && transactImpl != nil && asOfSnapshotImpl != nil {
		var base *BaseChannel
		if base, err = NewBaseChannel(
			ConnectionReferenceType,
			source, map[string]edn.Serializable{
				LabelReferenceProperty: label,
			}); err == nil {
			channel = &BaseConnectionChannel{
				BaseChannel:       base,
				transactImpl:      transactImpl,
				asOfSnapshotImpl:  asOfSnapshotImpl,
				tenantsWithSchema: make(map[string]bool),
			}
		}
	} else {
		err = edn.MakeError(edn.ErrInvalidInput, "label, transactor or snapshot are not valid")
	}

	return channel, err
}

// Transact the data to the channel
func (channel *BaseConnectionChannel) Transact(data ...interface{}) (result Result, err error) {

	var transactions []edn.Serializable
	if len(data) > 0 {
		for _, item := range data {

			switch typedItem := item.(type) {
			case string:
				transactions = append(transactions, RawString(typedItem))
			case edn.Serializable:
				transactions = append(transactions, typedItem)
			default:
				err = edn.MakeErrorWithFormat(edn.ErrInvalidInput, "Unsupported type: %T", typedItem)
			}
		}

	} else {
		err = edn.MakeError(edn.ErrInvalidInput, "No data")
	}

	if err == nil && len(transactions) > 0 {
		for _, trx := range transactions {
			result, err = channel.transactImpl(trx)
		}
	}

	return result, err
}

// Label to this particular channel
func (channel *BaseConnectionChannel) Label() string {
	return channel.BaseChannel.Label()
}

// AsOfSnapshot returns the snapshot channel as of the rules provided.
func (channel *BaseConnectionChannel) AsOfSnapshot(data interface{}) (snap SnapshotChannel, err error) {

	var elem edn.Serializable
	if data != nil {
		elem, err = decodeSerializable(data)
	}

	if err == nil {
		snap, err = channel.asOfSnapshotImpl(elem)
	}

	return snap, err
}

// LatestSnapshot returns the latest snapshot channel.
func (channel *BaseConnectionChannel) LatestSnapshot() (SnapshotChannel, error) {
	return channel.AsOfSnapshot(nil)
}
