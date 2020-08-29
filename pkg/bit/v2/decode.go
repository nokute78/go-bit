/*
   Copyright 2020 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package bit

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/goccy/go-reflect"
	"io"
	"strings"
)

const (
	tagKeyName = "bit"
)

var errCannotInterface = errors.New("CanInterface returns false")

// tagConfig represents StructTag.
//   "-"   : ignore the field
//   "skip": ignore but offset will be updated
//   "BE"  : the field is treated as big endian
//   "LE"  : the field is treated as little endian
type tagConfig struct {
	ignore bool
	skip   bool
	endian binary.ByteOrder
}

func parseStructTag(t reflect.StructTag) *tagConfig {
	s, ok := t.Lookup(tagKeyName)
	if !ok {
		return nil
	}
	ret := &tagConfig{}

	strs := strings.Split(s, ",")
	for _, v := range strs {
		switch v {
		case "-":
			ret.ignore = true
			return ret
		case "skip":
			ret.skip = true
			return ret
		case "BE":
			ret.endian = binary.BigEndian
		case "LE":
			ret.endian = binary.LittleEndian
		}

	}
	return ret
}

// Size returns size of v in bits.
/*
func Size(v interface{}) int {
	val := reflect.ValueOf(v)
	var i int = 0
	sizeOfValueInBits(&i, val, false)
	return i
}
*/

// This function will be panic if v doesn't support Bits function.
// if structtag is true, the function respects struct tag.
func sizeOfValueInBits(c *int, v reflect.Value, structtag bool) {
	switch v.Kind() {
	case reflect.Struct:
		if structtag {
			for i := 0; i < v.Type().NumField(); i++ {
				f := v.Type().Field(i)
				cnf := parseStructTag(f.Tag)
				if cnf != nil && cnf.ignore {
					continue
				}
				sizeOfValueInBits(c, v.Field(i), structtag)
			}
		} else {
			for i := 0; i < v.NumField(); i++ {
				sizeOfValueInBits(c, v.Field(i), structtag)
			}
		}
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			return
		}
		var elemSize int
		sizeOfValueInBits(&elemSize, v.Index(0), structtag)
		*c += (elemSize * v.Len())
	case reflect.Bool:
		*c += 1
	default:
		/* int, uint, float familiy */
		*c += v.Type().Bits()
	}
}

// filldata reads from b and fill v.
func fillData(b []byte, order binary.ByteOrder, v reflect.Value, o *Offset) error {
	var off Offset
	var err error
	var val reflect.Value

	if !v.CanInterface() {
		// skip unexported field
		var size int
		sizeOfValueInBits(&size, v, false)
		*o, err = o.AddOffset(Offset{Bit: uint64(size)})
		if err != nil {
			return err
		}
		return errCannotInterface
	}
	d := v.Interface()

	switch d.(type) {
	case uint8:
		ret, err := GetBitsAsByte(b, *o, 8, binary.LittleEndian)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(ret[0])
		off = Offset{1, 0}
	case uint16:
		ret, err := GetBitsAsByte(b, *o, 16, binary.LittleEndian)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(order.Uint16(ret))
		off = Offset{2, 0}
	case uint32:
		ret, err := GetBitsAsByte(b, *o, 32, binary.LittleEndian)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(order.Uint32(ret))
		off = Offset{4, 0}
	case uint64:
		ret, err := GetBitsAsByte(b, *o, 64, binary.LittleEndian)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(order.Uint64(ret))
		off = Offset{8, 0}

	case Bit:
		ret, err := GetBitsBitEndian(b, *o, 1, order)
		if err != nil {
			return err
		}
		if ret[0] {
			val = reflect.ValueOf(Bit(true))
		} else {
			val = reflect.ValueOf(Bit(false))
		}
		off = Offset{0, 1}
	default: /* other data types */
		switch v.Kind() {
		case reflect.Array:
			if v.Len() > 0 {
				if v.Index(0).Kind() == reflect.Bool {
					/* Bit array */
					/* when order is Big Endian, we should read the entire bits at once */
					ret, err := GetBitsBitEndian(b, *o, uint64(v.Len()), order)
					if err != nil {
						return err
					}
					// set data
					for i := 0; i < v.Len(); i++ {
						if v.Index(i).CanSet() {
							v.Index(i).Set(reflect.ValueOf(Bit(ret[i])))
						}
					}
					*o, err = o.AddOffset(Offset{Bit: uint64(v.Len())})
					if err != nil {
						return err
					}
					return nil
				} else if v.Index(0).Kind() == reflect.Uint8 {
					ret, err := GetBitsAsByte(b, *o, uint64(v.Len()*8), binary.LittleEndian)
					if err != nil {
						return err
					}
					for i := 0; i < v.Len(); i++ {
						if v.Index(i).CanSet() {
							if order == binary.BigEndian {
								// workaround! binary.Read doesn't support []byte in BigEndian
								v.Index(i).Set(reflect.ValueOf(ret[v.Len()-1-i]))
							} else {
								v.Index(i).Set(reflect.ValueOf(ret[i]))
							}
						}
					}
					*o, err = o.AddOffset(Offset{Byte: uint64(v.Len())})
					if err != nil {
						return err
					}
					return nil
				} else {
					for i := 0; i < v.Len(); i++ {
						err := fillData(b, order, v.Index(i), o)
						if err != nil && err != errCannotInterface {
							return err
						}
					}
					return nil
				}
			}
		case reflect.Struct:
			for i := 0; i < v.Type().NumField(); i++ {
				f := v.Type().Field(i)
				cnf := parseStructTag(f.Tag)
				if cnf != nil {
					/* struct tag is defined */
					if cnf.ignore {
						continue
					} else if cnf.skip {
						var bitSize int
						/* only updates offset. not fill. */
						sizeOfValueInBits(&bitSize, v.Field(i), true)
						*o, err = o.AddOffset(Offset{Bit: uint64(bitSize)})
						if err != nil {
							return err
						}
						continue
					} else if cnf.endian != nil {
						err := fillData(b, cnf.endian, v.Field(i), o)
						if err != nil && err != errCannotInterface {
							return err
						}
						continue
					}
				}
				err := fillData(b, order, v.Field(i), o)
				if err != nil && err != errCannotInterface {
					return err
				}
			}
			return nil
		default:
			return fmt.Errorf("Not Supported %s", v.Kind())
		}
	}

	// primitives
	if v.CanSet() {
		v.Set(val)
	} else {
		return fmt.Errorf("can not set %v\n", v)
	}
	*o, err = o.AddOffset(off)
	if err != nil {
		return err
	}

	return nil
}

// Read reads structured binary data from i into data.
// Data must be a pointer to a fixed-size value.
// Not exported struct field is ignored.
//   Supports StructTag.
//       `bit:"skip"` : ignore the field. Skip X bits which is the size of the field. It is useful for reserved field.
//       `bit:"-"`    : ignore the field. Offset is not changed.
func Read(r io.Reader, order binary.ByteOrder, data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		var c int = 0
		sizeOfValueInBits(&c, reflect.Indirect(v), true)
		byteSize := sizeOfBits(c)
		barr := make([]byte, byteSize)
		n, err := r.Read(barr)
		if err != nil {
			return err
		} else if n != byteSize {
			return fmt.Errorf("bit.Read:short read, expect=%d byte, read=%d byte", byteSize, n)
		}
		err = fillData(barr, order, reflect.Indirect(v), &Offset{})
		if err != io.EOF && err != errCannotInterface {
			return err
		}
	default:
		return binary.Read(r, order, data)
	}
	return nil
}
