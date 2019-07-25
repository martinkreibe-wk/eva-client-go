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

// Reference
type Reference interface {
	edn.CollectionElement

	// Type of reference
	Type() ChannelType

	// AddProperty
	AddProperty(name string, value edn.Element) error

	// GetProperty
	GetProperty(name string) (edn.Element, error)
}

const (
	ErrInvalidSerializer   = edn.ErrorMessage("Invalid serializer")
	LabelReferenceProperty = "label"
	AsOfReferenceProperty  = "as-of"
)

type refImpl struct {
	edn.CollectionElement
}

// newReference creates a new request.
func newReference(refType ChannelType, properties map[string]edn.Element) (Reference, error) {

	var refMap edn.CollectionElement
	var err error

	if refMap, err = edn.NewMap(); err != nil {
		return nil, err
	}

	ref := &refImpl{refMap}
	if err = ref.SetTag(string(refType)); err != nil {
		return nil, err
	}

	for name, value := range properties {
		if value != nil {
			err := ref.AddProperty(name, value)
			if err != nil {
				return nil, err
			}
		}
	}

	return ref, nil
}

// Type of this reference
func (ref *refImpl) Type() ChannelType {
	return ChannelType(ref.Tag())
}

// AddProperty will add the property by name, or if the value is nil, will remove it.
func (ref *refImpl) AddProperty(name string, value edn.Element) error {

	var symbol edn.SymbolElement
	var err error
	if symbol, err = edn.NewKeywordElement(name); err != nil {
		return err
	}

	return ref.Set(symbol, value)
}

// GetProperty returns the property by name
func (ref *refImpl) GetProperty(name string) (edn.Element, error) {
	var symbol edn.SymbolElement
	var err error
	if symbol, err = edn.NewKeywordElement(name); err != nil {
		return nil, err
	}

	return ref.Get(symbol)
}

func NewConnectionReference(label string) (Reference, error) {

	elem, err := edn.NewStringElement(label)
	if err != nil {
		return nil, err
	}

	return newReference(ConnectionReferenceType, map[string]edn.Element{
		LabelReferenceProperty: elem,
	})
}

func NewSnapshotAsOfReference(label string, asOf interface{}) (ref Reference, err error) {

	elem, err := edn.NewStringElement(label)
	if err != nil {
		return nil, err
	}

	properties := map[string]edn.Element{
		LabelReferenceProperty: elem,
	}

	if asOf != nil {
		var asOfElem edn.Element
		if asOfElem, err = edn.NewPrimitiveElement(asOf); err != nil {
			return nil, err
		}

		properties[AsOfReferenceProperty] = asOfElem
	}

	return newReference(SnapshotReferenceType, properties)
}

func NewSnapshotReference(label string) (req Reference, err error) {
	return NewSnapshotAsOfReference(label, nil)
}
