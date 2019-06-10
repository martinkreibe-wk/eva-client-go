# EDN processor

This package is the [EDN](https://github.com/edn-format/edn) processor for the Eva project. This library parses and serializes EDN.

## Architecture

The architecture of this package is as such, there are three main interfaces `edn.Element`, `edn.SymbolElement` and `edn.CollectionElement`
which house the actual values types. 


### Element interface

The `Element` interface is the wrapper for the raw data for an EDN supported type.

```go
// Element defines the interface for EDN elements.
type Element interface {

	// Serializer mixin
	Serializer

	// ElementType returns the current type of this element.
	ElementType() ElementType

	// Value of the element
	Value() interface{}

	// HasTag returns true if the element has a tag prefix
	HasTag() bool

	// Tag returns the prefixed tag if it exists.
	Tag() string

	// SetTag sets the tag to the incoming value. If the value is an empty string then the tag is unset.
	SetTag(string) (err error)

	// Equals checks if the input element is equal to this element.
	Equals(e Element) (result bool)
}
```

The `CollectionElement` interface is the wrapper for collection the raw data for an EDN supported type. Note that
internal structures are either an slice of `Element`s or a map of `Element`s.

```go
// CollectionElement defines the element for the EDN grouping construct. A group is a sequence of values. Groups are
// represented by zero or more elements enclosed in parentheses (). Note that Group can be heterogeneous.
type CollectionElement interface {
	Element

	// Len return the quantity of items in this collection.
	Len() int

	// IterateChildren will iterator over all children.
	IterateChildren(iterator ChildIterator) (err error)

	// Append the elements into this collection.
	Append(children ...Element) (err error)

	// Get the key from the collection.
	Get(key interface{}) (Element, error)

	// Merge one collection into another.
	Merge(CollectionElement) error
}
```

The `SymbolElement` interface handles the special case for keywords and symbols which need special rules to handle
prefixes.

```go
// Symbols are used to represent identifiers, and should map to something other than strings, if possible.
type SymbolElement interface {
	Element

	// Modifier for this symbol
	Modifier() string

	// Prefix to this symbol
	Prefix() string

	// Name to this symbol
	Name() string

	// AppendNameOntoNamespace will append the input name onto the namespace.
	AppendNameOntoNamespace(string) string
}
```
 
The following table describes the supported types and their affiliations in various
contexts.

| EDN Type | Golang Type | Eva                                     |
|----------|-------------|-----------------------------------------|
| [boolean](https://github.com/edn-format/edn#booleans) | `bool` | `db.type\boolean` |
| [character](https://github.com/edn-format/edn#characters) | `rune` | `db.type\character` |
| [float](https://github.com/edn-format/edn#floating-point-numbers) | `float64` | `db.type\float` |
| [instant](https://github.com/edn-format/edn#inst-rfc-3339-format) | `time.Time` | `db.type\instant` |
| [integer](https://github.com/edn-format/edn#integers) | `int64` | `db.type\long` |
| [keyword](https://github.com/edn-format/edn#keywords) | `edn.SymbolElement` | `db.type\keyword` |
| [nil](https://github.com/edn-format/edn#nil) | `interface{}` set to `nil` | `db.type\nil` |
| [string](https://github.com/edn-format/edn#strings) |  `string` | `db.type\string` |
| [symbol](https://github.com/edn-format/edn#symbols) | `edn.SymbolElement` | `db.type\symbol` |
| [UUID](https://github.com/edn-format/edn#uuid-f81d4fae-7dec-11d0-a765-00a0c91e6bf6) |  `github.com/mattrobenolt/gocql/uuid.UUID` |  `db.type\uuid` |
| [list](https://github.com/edn-format/edn#lists) | `edn.CollectionElement` | `db.type\group` |
| [map](https://github.com/edn-format/edn#maps) | `edn.CollectionElement` | `db.type\map` |
| [set](https://github.com/edn-format/edn#sets) | `edn.CollectionElement` | `db.type\set` |
| [vector](https://github.com/edn-format/edn#vectors) | `edn.CollectionElement` |  `db.type\vector` |


The following are types that are known Eva types that are yet to be supported:

* uri
* bytes
* bigint
* bigdec
* double
* ref

## Usage

### Parsing

There are two functions used to parse the edn into the appropriate elements:
1. `Parse(string) (Element, error)`
    Parse the string into an element. The result can be either an `Element` or `CollectionElement` with any errors
    reporting through the `error` return parameter.
    
1. `ParseCollection(string) (CollectionElement, error)`
    Parse the string assuming the result of the parse will be a `CollectionElement` and reporting any issues through
    the `error` return parameter.

### Generating primitive elements

Use `NewPrimitiveElement(interface{}) (Element, error)` to generate a primitive from any supported type. Otherwise
specific `Element`s can be created by using the specific constructors:

* `NewBooleanElement(bool) (Element)`
* `NewCharacterElement(rune) (Element)`
* `NewFloatElement(float64) (Element)`
* `NewInstantElement(time.Time) (Element)`
* `NewIntegerElement(int64) (Element)`
* `NewKeywordElement(...string) (SymbolElement, error)`
* `NewList(...Element) (CollectionElement, error)`
* `NewMap(...Pair) (CollectionElement, error)`
* `NewNilElement() (Element)`
* `NewSet(...Element) (CollectionElement, error)`
* `NewStringElement(string) (Element)`
* `NewSymbolElement(...string) (SymbolElement, error)`
* `NewUUIDElement(uuid.UUID) (Element)`
* `NewVector(...Element) (CollectionElement, error)`

## Testing

This package uses [ginkgo](https://onsi.github.io/ginkgo/) and [gomega](https://onsi.github.io/gomega/) to facilitate
testing. There are a large number of tests due to the nature of parsing (for full tests, >20k tests).

For the general tests, just run: `go test github.com/Workiva/eva-client-go/edn`

Add the following for more details:

| Test Configuration                       | Description                                                          |
|------------------------------------------|----------------------------------------------------------------------|
| `-cover` command line flag               | shows the code coverage from the tests.                              |
| `-test.coverprofile=<file>` command line | saves the coverage profile to the `<file> location                   |
| `-test.v` command line flag              | show all verbose details on execution.                               |
| `FULL_TESTS=true` environment variable   | runs tests that are by default not on. These are very verbose (~20x) |

Here is the example for running all tests with code coverage and saving coverage details to the 'coverage.edn.out' file.
```FULL_TESTS=true go test -cover -test.v -test.coverprofile=coverage.edn.out github.com/Workiva/eva-client-go/edn```
