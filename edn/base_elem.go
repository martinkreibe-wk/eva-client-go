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

package edn

import (
	"strings"
)

// baseElement defines the base element features.
type baseElemImpl struct {

	// elemType is the type this element houses.
	elemType ElementType

	// equality is the tester for equality
	equality equalityFunc

	// tag of this element.
	tag string

	// value of this element.
	value interface{}
}

// Equals checks if the input element is equal to this element.
func (elem *baseElemImpl) Equals(e Element) (result bool) {
	if elem.ElementType() == e.ElementType() {
		if elem.Tag() == e.Tag() {
			result = elem.equality(elem, e)
		}
	}
	return result
}

// ElementType returns the current type of this element.
func (elem *baseElemImpl) ElementType() ElementType {
	return elem.elemType
}

// HasTag returns true if the element has a tag prefix
func (elem *baseElemImpl) HasTag() bool {
	return len(elem.tag) != 0
}

// Tag returns the prefixed tag if it exists.
func (elem *baseElemImpl) Tag() string {
	return elem.tag
}

// tagged elements
//
// edn supports extensibility through a simple mechanism. # followed immediately by a symbol starting with an alphabetic
// character indicates that that symbol is a tag. A tag indicates the semantic interpretation of the following element.
// It is envisioned that a reader implementation will allow clients to register handlers for specific tags. Upon
// encountering a tag, the reader will first read the next element (which may itself be or comprise other tagged
// elements), then pass the result to the corresponding handler for further interpretation, and the result of the
// handler will be the data value yielded by the tag + tagged element, i.e. reading a tag and tagged element yields one
// value. This value is the value to be returned to the program and is not further interpreted as edn data by the
// reader.
//
// This process will bottom out on elements either understood or built-in.
//
// Thus you can build new distinct readable elements out of (and only out of) other readable elements, keeping extenders
// and extension consumers out of the text business.
//
// The semantics of a tag, and the type and interpretation of the tagged element are defined by the steward of the tag.
//
// #myapp/Person {:first "Fred" :last "Mertz"}
//
// If a reader encounters a tag for which no handler is registered, the implementation can either report an error, call
// a designated 'unknown element' handler, or create a well-known generic representation that contains both the tag and
// the tagged element, as it sees fit. Note that the non-error strategies allow for readers which are capable of reading
// any and all edn, in spite of being unaware of the details of any extensions present.

// rules for tags
//
// Tag symbols without a prefix are reserved by edn for built-ins defined using the tag system.
// User tags must contain a prefix component, which must be owned by the user (e.g. trademark or domain) or known unique
// in the communication context.
// A tag may specify more than one format for the tagged element, e.g. both a string and a vector representation.
// Tags themselves are not elements. It is an error to have a tag without a corresponding tagged element.

// SetTag sets the tag to the incoming value. If the value is an empty string then the tag is unset.
func (elem *baseElemImpl) SetTag(value string) (err error) {

	if len(value) > 0 {
		tag := value
		if strings.HasPrefix(value, TagPrefix) {
			tag = strings.TrimPrefix(value, TagPrefix)
		}

		var prefix, name string
		if prefix, name, err = decodeSymbol(tag); err == nil {
			elem.tag = encodeSymbol(prefix, name)
		}
	} else {
		elem.tag = ""
	}

	return err
}

// Value return the raw representation of this element.
func (elem *baseElemImpl) Value() interface{} {
	return elem.value
}
