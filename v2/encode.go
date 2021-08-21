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
	"fmt"
	"github.com/goccy/go-reflect"
	"io"
)

// write writes v to b.
func write(v reflect.Value, order binary.ByteOrder, b []byte, o *Offset) error {
	var off Offset
	var err error

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
		val := d.(uint8)
		bits, err := GetBits([]byte{byte(val)}, Offset{0, 0}, 8, order)
		if err != nil {
			return err
		}
		if err := SetBits(b, *o, bits, binary.LittleEndian); err != nil {
			return err
		}
		off = Offset{1, 0}
	case uint16:
		val := d.(uint16)
		bs := make([]byte, 2)
		order.PutUint16(bs, val)
		bits, err := GetBits(bs, Offset{0, 0}, 16, binary.LittleEndian)
		if err != nil {
			return err
		}
		if err := SetBits(b, *o, bits, binary.LittleEndian); err != nil {
			return err
		}
		off = Offset{2, 0}

	case uint32:
		val := d.(uint32)
		bs := make([]byte, 4)
		order.PutUint32(bs, val)
		bits, err := GetBits(bs, Offset{0, 0}, 32, binary.LittleEndian)
		if err != nil {
			return err
		}
		if err := SetBits(b, *o, bits, binary.LittleEndian); err != nil {
			return err
		}
		off = Offset{4, 0}
	case uint64:
		val := d.(uint64)
		bs := make([]byte, 8)
		order.PutUint64(bs, val)
		bits, err := GetBits(bs, Offset{0, 0}, 64, binary.LittleEndian)
		if err != nil {
			return err
		}
		if err := SetBits(b, *o, bits, binary.LittleEndian); err != nil {
			return err
		}
		off = Offset{8, 0}
	case Bit:
		val := d.(Bit)
		if err := SetBitsBitEndian(b, *o, []Bit{val}, order); err != nil {
			return err
		}
		off = Offset{0, 1}
	default:
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			if v.Len() > 0 {
				if v.Index(0).Kind() == reflect.Bool {
					/* Bit array */
					/* when order is Big Endian, we should read the entire bits at once */
					bitSize := v.Len()
					bits := make([]Bit, bitSize)
					for i := 0; i < bitSize; i++ {
						bits[i] = Bit(v.Index(i).Bool())
					}

					if err := SetBitsBitEndian(b, *o, bits, order); err != nil {
						return err
					}
					off = Offset{0, uint64(bitSize)}

				} else if v.Index(0).Kind() == reflect.Uint8 {
					// byte slice / byte array
					var bs []byte
					if v.Kind() == reflect.Slice {
						bs = v.Bytes()
					} else {
						// array
						bs = make([]byte, v.Len())
						for i := 0; i < v.Len(); i++ {
							bs[i] = byte(v.Index(i).Uint())
						}
					}
					bits, err := GetBits(bs, Offset{0, 0}, uint64(v.Len()*8), order)
					if err != nil {
						return err
					}
					if err := SetBits(b, *o, bits, binary.LittleEndian); err != nil {
						return err
					}
					off = Offset{uint64(v.Len()), 0}
				} else {
					for i := 0; i < v.Len(); i++ {
						err := write(v.Index(i), order, b, o)
						if err != nil && err != errCannotInterface {
							return err
						}
					}
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
						err := write(v.Field(i), cnf.endian, b, o)
						if err != nil && err != errCannotInterface {
							return err
						}
						continue
					}
				}
				err := write(v.Field(i), order, b, o)
				if err != nil && err != errCannotInterface {
					return err
				}
			}
			return nil
		default:
			return fmt.Errorf("Not Supported %s", v.Kind())
		}
	}

	*o, err = o.AddOffset(off)

	return err
}

// Write writes structured binary data from input into w.
func Write(w io.Writer, order binary.ByteOrder, input interface{}) error {
	v := reflect.ValueOf(input)
	var vv reflect.Value

	switch v.Kind() {
	case reflect.Ptr:
		vv = reflect.Indirect(reflect.ValueOf(input))
	case reflect.Array, reflect.Slice, reflect.Struct:
		vv = reflect.ValueOf(input)
	default:
		return binary.Write(w, order, input)
	}

	var c int = 0
	var off Offset
	sizeOfValueInBits(&c, vv, true)
	byteSize := sizeOfBits(c)
	barr := make([]byte, byteSize)

	err := write(vv, order, barr, &off)
	_, err = w.Write(barr)
	return err
}
