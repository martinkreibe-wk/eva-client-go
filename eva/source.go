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

	// ErrDuplicateSourceType defines the duplicate source type error.
	ErrDuplicateSourceType = edn.ErrorMessage("Duplicate source type added")

	// ErrUnknownSourceType defines the unknown source type error.
	ErrUnknownSourceType = edn.ErrorMessage("Unknown source type")
)

// Source of eva information.
type Source interface {

	// Serializer returns the serializer to use.
	Serializer() (edn.Serializer, error)

	// Tenant gets the current tenant.
	Tenant() Tenant

	// Connection returns the connection channel
	Connection(label interface{}) (conn ConnectionChannel, err error)

	// LatestSnapshot returns the latest snapshot channel for the specified category of data.
	LatestSnapshot(label interface{}) (SnapshotChannel, error)

	// AsOfSnapshot returns the snapshot channel of the specified category as of the rules provided.
	AsOfSnapshot(label interface{}, asOf interface{}) (SnapshotChannel, error)

	// Query the source for data.
	Query(query interface{}, parameters ...interface{}) (Result, error)
}

// sourceFactory defines the mechanism for creating a source.
type sourceFactory func(config Configuration, tenant Tenant) (source Source, err error)

// sourceFactories holds the various sources.
var sourceFactories = map[string]sourceFactory{}

// AddSourceFactory will add source factory to the collection of factories.
func AddSourceFactory(sourceType string, factory sourceFactory) (err error) {
	if _, has := sourceFactories[sourceType]; !has {
		sourceFactories[sourceType] = factory
	} else {
		err = edn.MakeError(ErrDuplicateSourceType, sourceType)
	}
	return err
}

// NewSource create a new source from the configuration.
func NewSource(config Configuration, tenant Tenant) (source Source, err error) {
	if config != nil {
		var srcConfig SourceConfiguration
		if srcConfig, err = config.Source(); err == nil {
			if factory, has := sourceFactories[srcConfig.Type()]; has {
				if srcConfig, err = config.Source(); err == nil {
					source, err = factory(config, tenant)
				}
			} else {
				err = edn.MakeError(ErrUnknownSourceType, srcConfig.Type())
			}
		}
	} else {
		err = edn.MakeError(ErrInvalidConfiguration, nil)
	}

	return source, err
}

// Sources returns the collection of sources
func Sources() (sources []string) {
	for s := range sourceFactories {
		sources = append(sources, s)
	}

	return sources
}

// HasSource checks if the source type exists.
func HasSource(source string) (has bool) {
	_, has = sourceFactories[source]
	return has
}
