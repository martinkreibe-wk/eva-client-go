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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattrobenolt/gocql/uuid"
)

const (

	// EvaEdnMimeType defines the mime type for the eva edn data.
	EvaEdnMimeType SerializerMimeType = "application/vnd.eva+edn"
)

type SerializerMimeType string

// String this Mime Type.
func (smt SerializerMimeType) String() string {
	return string(smt)
}

func (smt SerializerMimeType) SerializeTo(stream Stream, elem Element) error {

	if stream == nil {
		return MakeError(ErrInvalidInput, "nil stream")
	}

	if elem == nil {
		return MakeError(ErrInvalidInput, "nil element")
	}

	if elem.HasTag() {
		tag := elem.Tag()
		if len(tag) > 0 {
			if _, err := stream.Write([]byte(TagPrefix + tag + " ")); err != nil {
				return err
			}
		}
	}

	value := elem.Value()
	switch t := elem.ElementType(); t {
	case BooleanType:
		_, err := stream.Write([]byte(strconv.FormatBool(value.(bool))))
		return err
	case CharacterType:
		val := value.(rune)

		// look at the special characters first.
		for v, r := range specialCharacters {
			if val == r {
				_, err := stream.Write([]byte(CharacterPrefix + v))
				return err
			}
		}

		// if there is no special character, then quote the rune, remove the single quotes around this, then
		// if it is an ASCII then make sure to prefix is intact.
		char := strings.Trim(fmt.Sprintf("%+q", val), "'")
		if strings.HasPrefix(char, CharacterPrefix) {
			_, err := stream.Write([]byte(char))
			return err
		}

		_, err := stream.Write([]byte(CharacterPrefix + char))
		return err
	case FloatType:
		_, err := stream.Write([]byte(strconv.FormatFloat(value.(float64), 'E', -1, 64)))
		return err
	case InstantType:
		_, err := stream.Write([]byte(value.(time.Time).Format(time.RFC3339)))
		return err
	case IntegerType:
		_, err := stream.Write([]byte(strconv.FormatInt(value.(int64), 10)))
		return err
	case NilType:
		_, err := stream.Write([]byte("nil"))
		return err
	case StringType:
		_, err := stream.Write([]byte(strconv.Quote(value.(string))))
		return err
	case SymbolType, KeywordType:
		symbol, is := value.(SymbolElement)
		if !is {
			return MakeErrorWithFormat(ErrInvalidElement, "Expected Symbol or Keyword, Got: %s", t)
		}
		_, err := stream.Write([]byte(symbol.AppendNameOntoNamespace(symbol.Name())))
		return err
	case UUIDType:
		_, err := stream.Write([]byte(value.(uuid.UUID).String()))
		return err
	case ListType, SetType, VectorType:
		val := value.(*collectionElemImpl)

		_, err := stream.Write([]byte(val.startSymbol))
		if err != nil {
			return err
		}

		first := true
		err = val.IterateChildren(func(_ Element, child Element) error {
			if first {
				first = false
			} else {
				_, e := stream.Write([]byte(val.separatorSymbol))
				if e != nil {
					return e
				}
			}

			if child != nil {
				if e := smt.SerializeTo(stream, child); e != nil {
					return e
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		_, err = stream.Write([]byte(val.endSymbol))
		return err
	case MapType:
		val := value.(*collectionElemImpl)

		_, err := stream.Write([]byte(val.startSymbol))
		if err != nil {
			return err
		}

		first := true
		err = val.IterateChildren(func(key Element, child Element) error {
			if first {
				first = false
			} else {
				_, e := stream.Write([]byte(val.separatorSymbol))
				if e != nil {
					return e
				}
			}

			if e := smt.SerializeTo(stream, key); e != nil {
				return e
			}
			_, e := stream.Write([]byte(val.keyValueSeparatorSymbol))
			if e != nil {
				return e
			}

			if child != nil {
				if e := smt.SerializeTo(stream, child); e != nil {
					return e
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		_, err = stream.Write([]byte(val.endSymbol))
		return err
	default:
		return MakeErrorWithFormat(ErrInvalidElement, "Unknown type: %s", t)
	}
}
