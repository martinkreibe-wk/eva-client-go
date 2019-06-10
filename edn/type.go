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

// ElementType indicates the EDN element construct
type ElementType string

// ElementTypeFactory defines the factory for an element.
type ElementTypeFactory func(interface{}) (Element, error)

const (

	// typeNamespace hold the type namespace
	typeNamespace = KeywordPrefix + "db.type"

	// UnknownType TODO
	UnknownType = ElementType("")

	// bytes, uri, ref
	// KeywordType is the value type for keywords. Keywords are used as names, and are interned for efficiency.
	// Keywords map to the native interned-name type in languages that support them.
	NilType = ElementType(typeNamespace + SymbolSeparator + "nil")

	// CharacterType TODO
	CharacterType = ElementType(typeNamespace + SymbolSeparator + "character")

	// KeywordType is the value type for keywords. Keywords are used as names, and are interned for efficiency.
	// Keywords map to the native interned-name type in languages that support them.
	KeywordType = ElementType(typeNamespace + SymbolSeparator + "keyword")

	// SymbolType TODO
	SymbolType = ElementType(typeNamespace + SymbolSeparator + "symbol")

	// StringType is the value type for strings.
	StringType = ElementType(typeNamespace + SymbolSeparator + "string")

	// BooleanType value type.
	BooleanType = ElementType(typeNamespace + SymbolSeparator + "boolean")

	// IntegerType is the fixed integer value type. Same semantics as a Java long: 64 bits wide, two's complement binary
	// representation.
	IntegerType = ElementType(typeNamespace + SymbolSeparator + "long")

	// BigIntType is the value type for arbitrary precision integers. Maps to java.math.BigInteger on Java platforms.
	BigIntType = ElementType(typeNamespace + SymbolSeparator + "bigint")

	// FloatType is the floating point value type. Same semantics as a Java float: single-precision 32-bit IEEE 754
	// floating point.
	FloatType = ElementType(typeNamespace + SymbolSeparator + "float")

	// DoubleType is the floating point value type. Same semantics as a Java double: double-precision 64-bit IEEE 754
	// floating point.
	DoubleType = ElementType(typeNamespace + SymbolSeparator + "double")

	// BigDecType is the value type for arbitrary precision floating point numbers. Maps to java.math.BigDecimal on Java
	// platforms.
	BigDecType = ElementType(typeNamespace + SymbolSeparator + "bigdec")

	// RefType is the value type for references. All references from one entity to another are through attributes with
	// this value type.
	RefType = ElementType(typeNamespace + SymbolSeparator + "ref")

	// TODO
	ListType   = ElementType(typeNamespace + SymbolSeparator + "group")
	VectorType = ElementType(typeNamespace + SymbolSeparator + "vector")
	MapType    = ElementType(typeNamespace + SymbolSeparator + "map")
	SetType    = ElementType(typeNamespace + SymbolSeparator + "set")

	// InstantType is the value type for instants in time. Stored internally as a number of milliseconds since midnight,
	// January 1, 1970 UTC. Maps to java.util.Date on Java platforms.
	InstantType = ElementType(typeNamespace + SymbolSeparator + "instant")

	// UUIDType is the value type for UUIDs. Maps to java.util.UUID on Java platforms.
	UUIDType = ElementType(typeNamespace + SymbolSeparator + "uuid")

	// URIType is the value type for URIs. Maps to java.net.URI on Java platforms.
	URIType = ElementType(typeNamespace + SymbolSeparator + "uri")

	// BytesType is the value type for small binary data. Maps to byte array on Java platforms. See limitations.
	BytesType = ElementType(typeNamespace + SymbolSeparator + "bytes")

	// ErrInvalidFactory defines the factory error
	ErrInvalidFactory = ErrorMessage("Invalid factory")

	// ErrInvalidInput defines the input error
	ErrInvalidInput = ErrorMessage("Invalid input")
)

// typeFactories hold the collection of element factories.
var typeFactories = map[ElementType]ElementTypeFactory{}

type elementDefinition struct {
	elemType    ElementType
	initializer func(lexer Lexer) error
}

var unknownTypeDef = &elementDefinition{UnknownType, nil}

// typeDefinitions holds the type to name/initializer mappings
// NOTE: ORDER MATTERS!!
var typeDefinitions = []*elementDefinition{
	unknownTypeDef,
	{NilType, initNil},
	{BooleanType, initBoolean},
	{StringType, initString},
	{CharacterType, initCharacter},
	{SymbolType, initSymbol},
	{KeywordType, initKeyword},
	{IntegerType, initInteger},
	{FloatType, initFloat},
	{InstantType, initInstant},
	{UUIDType, initUUID},
	{ListType, initList},
	{VectorType, initVector},
	{MapType, initMap},
	{SetType, initSet},

	// TODO
	{URIType, nil},
	{BytesType, nil},
	{BigIntType, nil},
	{BigDecType, nil},
	{DoubleType, nil},
	{RefType, nil},
}

// init will initialize the package - NOTE this is not testable
func init() {
	initAll()
}

// addElementTypeFactory adds an element factory to the factory collection.
func addElementTypeFactory(elemType ElementType, elemFactory ElementTypeFactory) (err error) {
	if _, has := typeFactories[elemType]; !has {
		typeFactories[elemType] = elemFactory
	} else {
		err = MakeError(ErrInvalidFactory, elemType)
	}

	return err
}

// init the package
func initAll() {

	var err error

	var lexer Lexer
	if lexer, err = getLexer(); err == nil {
		for _, def := range typeDefinitions {
			if def.initializer != nil {
				if err = def.initializer(lexer); err != nil {
					break
				}
			}
		}
	}

	if err != nil {
		panic(err)
	}
}

// IsCollection indicates that this type is a collection
func (t ElementType) IsCollection() (isColl bool) {
	switch t {
	case ListType, VectorType, MapType, SetType:
		isColl = true
	}
	return isColl
}

// Name returns the name of the set.
func (t ElementType) Name() (name string) {
	return string(t)
}
