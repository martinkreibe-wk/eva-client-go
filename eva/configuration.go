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
	"encoding/json"
	"github.com/Workiva/eva-client-go/edn"
)

const (

	// ErrInvalidConfiguration defines an invalid configuration.
	ErrInvalidConfiguration = edn.ErrorMessage("Invalid configuration")
)

// SourceConfiguration is the configuration for a source.
type SourceConfiguration interface {

	// Type of source.
	Type() string

	// Serializer this source emanates with.
	Serializer() (edn.Serializer, error)

	// Setting within this part of the configuration.
	Setting(name string) (string, bool)
}

// Configuration for an eva connection.
type Configuration interface {

	// Category defined in the configuration
	Category() string

	// Source as defined in the configuration
	Source() (SourceConfiguration, error)
}

// sourceConfigImpl defines the configuration for the source.
type sourceConfigImpl map[string]string

// Type of source.
func (config sourceConfigImpl) Type() string {
	return config["type"]
}

// Setting within this part of the configuration.
func (config sourceConfigImpl) Setting(name string) (value string, has bool) {
	value, has = config[name]
	return value, has
}

// Serializer this source emanates with.
func (config sourceConfigImpl) Serializer() (serializer edn.Serializer, err error) {

	serializer = edn.EvaEdnMimeType
	if value, has := config["mime"]; has {
		serializer = edn.SerializerMimeType(value)
	}

	return serializer, err
}

// configImpl implements the source configuration.
type configImpl struct {

	// CategoryValue defines the category for this configuration.
	CategoryValue string `json:"category"`

	// ServiceName is the name of the service.
	SourceData sourceConfigImpl `json:"source"`
}

// Category defined in the configuration
func (config *configImpl) Category() string {
	return config.CategoryValue
}

// Source as defined in the configuration
func (config *configImpl) Source() (sc SourceConfiguration, err error) {
	return config.SourceData, err
}

// NewConfiguration creates a new configuration from the input string.
func NewConfiguration(data string) (config Configuration, err error) {

	config = &configImpl{}
	if err = json.Unmarshal([]byte(data), config); err != nil {
		config = nil
	}

	return config, err
}
