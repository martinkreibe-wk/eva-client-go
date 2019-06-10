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

// Tenant defines the tenant object.
type Tenant interface {

	// Name of the tenant
	Name() string

	// CorrelationId returns the correlation id if it is set. If the correlation id is not set, this will return false.
	CorrelationId() (string, bool)
}

// NewTenant creates a new tenant.
func NewTenant(name string) (Tenant, error) {
	return NewCorrelationTenant(name, "")
}

// NewCorrelationTenant creates a new tenant with a correlation id.
func NewCorrelationTenant(name string, correlation string) (t Tenant, err error) {
	return &tenantImpl{
		name:        name,
		correlation: correlation,
	}, nil
}
