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

type ConnectionChannelMaker func(label edn.Element, source Source) (channel ConnectionChannel, err error)

// QueryImplementation defines the query implementation function.
type QueryImplementation func(interface{}, ...interface{}) (Result, error)

// BaseSource defines the base source.
type BaseSource struct {
	source Source
	config Configuration
	tenant Tenant
	maker  ConnectionChannelMaker
	query  QueryImplementation
}

// NewBaseSource creates a new source query.
func NewBaseSource(config Configuration, tenant Tenant, srcImpl Source, maker ConnectionChannelMaker, impl QueryImplementation) (source *BaseSource, err error) {

	if impl != nil && maker != nil {
		if config != nil {
			if len(config.Category()) > 0 {
				if _, err = config.Source(); err == nil {
					source = &BaseSource{
						source: srcImpl,
						config: config,
						tenant: tenant,
						maker:  maker,
						query:  impl,
					}
				}
			} else {
				err = edn.MakeError(ErrInvalidConfiguration, "Invalid category")
			}
		} else {
			err = edn.MakeError(ErrInvalidConfiguration, nil)
		}
	} else {
		err = edn.MakeError(ErrInvalidConfiguration, "Source query implementation or connection finder is nil")
	}

	return source, err
}

func (source *BaseSource) Tenant() Tenant {
	return source.tenant
}

func (source *BaseSource) Category() string {
	return source.config.Category()
}

func (source *BaseSource) Connection(label interface{}) (ConnectionChannel, error) {

	var elem edn.Element
	switch v := label.(type) {
	case edn.Element:
		elem = v
	default:
		var err error
		elem, err = edn.NewPrimitiveElement(label)
		if err != nil {
			return nil, err
		}
	}

	return source.maker(elem, source.source)
}

// Serializer returns the serializer this source has.
func (source *BaseSource) Serializer() (serializer edn.Serializer, err error) {
	var src SourceConfiguration
	if src, err = source.config.Source(); err == nil {
		serializer, err = src.Serializer()
	}
	return serializer, err
}

// LatestSnapshot returns the latest snapshot channel for the specified category of data.
func (source *BaseSource) LatestSnapshot(label interface{}) (channel SnapshotChannel, err error) {
	var conn ConnectionChannel
	if conn, err = source.Connection(label); err == nil {
		channel, err = conn.LatestSnapshot()
	}

	return channel, err
}

// AsOfSnapshot returns the snapshot channel of the specified category as of the rules provided.
func (source *BaseSource) AsOfSnapshot(label interface{}, asOf interface{}) (channel SnapshotChannel, err error) {
	var conn ConnectionChannel
	if conn, err = source.Connection(label); err == nil {
		channel, err = conn.AsOfSnapshot(asOf)
	}

	return channel, err
}

// Query the source for data.
func (source *BaseSource) Query(query interface{}, parameters ...interface{}) (result Result, err error) {
	return source.query(query, parameters...)
}
