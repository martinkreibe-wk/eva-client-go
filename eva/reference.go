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
	edn.Serializable

	// Type of reference
	Type() ChannelType

	// AddProperty
	AddProperty(name string, value edn.Serializable) error

	// GetProperty
	GetProperty(name string) edn.Serializable
}

const (
	ErrInvalidSerializer   = edn.ErrorMessage("Invalid serializer")
	LabelReferenceProperty = "label"
	AsOfReferenceProperty  = "as-of"
)

type refImpl struct {
	refType    ChannelType
	properties map[string]edn.Serializable
}

// newReference creates a new request.
func newReference(refType ChannelType, properties map[string]edn.Serializable) (ref Reference, err error) {

	if properties == nil {
		properties = make(map[string]edn.Serializable)
	}

	ref = &refImpl{
		refType:    refType,
		properties: properties,
	}

	return ref, err
}

func (ref *refImpl) String() string {
	var str string
	PanicOnError(func() error {
		var err error
		str, err = ref.Serialize(edn.EvaEdnMimeType)
		return err
	})
	return str
}

// Serialize the reference.
func (ref *refImpl) Serialize(with edn.Serializer) (value string, err error) {
	if with != nil {
		var elem edn.CollectionElement
		if elem, err = edn.NewMap(); err == nil {
			for name, value := range ref.properties {
				if value != nil {
					var symbol edn.SymbolElement
					if symbol, err = edn.NewKeywordElement(name); err == nil {
						switch v := value.(type) {
						case edn.Element:
							err = elem.Append(symbol, v)
						case rawStringImpl:
							err = elem.Append(symbol, edn.NewStringElement(v.String()))
						case rawIntImpl:
							err = elem.Append(symbol, edn.NewIntegerElement(v.Int()))
						default:
							err = edn.MakeError(edn.ErrInvalidInput, value)
						}
					}
				}

				if err != nil {
					break
				}
			}
		}

		if err == nil {
			if err = elem.SetTag(string(ref.Type())); err == nil {
				value, err = elem.Serialize(with)
			}
		}
	} else {
		err = edn.MakeError(ErrInvalidSerializer, nil)
	}

	return value, err
}

// Type of this reference
func (ref *refImpl) Type() ChannelType {
	return ref.refType
}

// AddProperty will add the property by name, or if the value is nil, will remove it.
func (ref *refImpl) AddProperty(name string, value edn.Serializable) error {
	if value != nil {
		ref.properties[name] = value
	} else {
		delete(ref.properties, name)
	}

	return nil
}

// GetProperty returns the property by name
func (ref *refImpl) GetProperty(name string) edn.Serializable {
	return ref.properties[name]
}

func NewConnectionReference(label string) (ref Reference, err error) {
	return newReference(ConnectionReferenceType, map[string]edn.Serializable{
		LabelReferenceProperty: RawString(label),
	})
}

func NewSnapshotAsOfReference(label string, asOf interface{}) (ref Reference, err error) {

	var asOfElem edn.Serializable
	if asOfElem, err = decodeSerializable(asOf); err == nil {
		ref, err = newReference(SnapshotReferenceType, map[string]edn.Serializable{
			LabelReferenceProperty: RawString(label),
			AsOfReferenceProperty:  asOfElem,
		})
	}

	return ref, err
}

func NewSnapshotReference(label string) (req Reference, err error) {
	return NewSnapshotAsOfReference(label, nil)
}
