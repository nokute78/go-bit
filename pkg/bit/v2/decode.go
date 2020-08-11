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
	"io"
	//	"reflect"
	"github.com/goccy/go-reflect"
)

// Size returns size of v in bits.
func Size(v interface{}) int {
	val := reflect.ValueOf(v)
	var i int = 0
	sizeOfValueInBits(&i, val)
	return i
}

// This function will be panic if v doesn't support Bits function.
func sizeOfValueInBits(c *int, v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			sizeOfValueInBits(c, f)
		}
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			return
		}
		i := v.Index(0)
		var elemSize int
		sizeOfValueInBits(&elemSize, i)
		*c += (elemSize * v.Len())
	case reflect.Bool:
		*c += 1
	default:
		/* int, uint, float familiy */
		*c += v.Type().Bits()
	}
}

func fillData(b []byte, order binary.ByteOrder, v reflect.Value, o *Offset) error {
	var off Offset
	var err error
	var val reflect.Value

	d := v.Interface()

	switch d.(type) {
	case uint8:
		ret, err := GetBitsAsByte(b, *o, 8)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(ret[0])
		off = Offset{1, 0}
	case uint16:
		ret, err := GetBitsAsByte(b, *o, 16)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(order.Uint16(ret))
		off = Offset{2, 0}
	case uint32:
		ret, err := GetBitsAsByte(b, *o, 32)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(order.Uint32(ret))
		off = Offset{4, 0}
	case uint64:
		ret, err := GetBitsAsByte(b, *o, 64)
		if err != nil {
			return err
		}
		val = reflect.ValueOf(order.Uint64(ret))
		off = Offset{8, 0}

	case Bit:
		ret, err := GetBitsAsByte(b, *o, 1)
		if err != nil {
			return err
		}
		if ret[0] > 0 {
			val = reflect.ValueOf(Bit(true))
		} else {
			val = reflect.ValueOf(Bit(false))
		}

		off = Offset{0, 1}
	default:
		switch v.Kind() {
		case reflect.Array:
			for i := 0; i < v.Len(); i++ {
				err := fillData(b, order, v.Index(i), o)
				if err != nil {
					return err
				}
			}
			return nil
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				err := fillData(b, order, v.Field(i), o)
				if err != nil {
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
		fmt.Println("can not set")
	}
	*o, err = o.AddOffset(off)
	if err != nil {
		return err
	}

	return nil
}

func Read(r io.Reader, order binary.ByteOrder, data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		var c int = 0
		sizeOfValueInBits(&c, reflect.Indirect(v))
		byteSize := sizeOfBits(c)
		barr := make([]byte, byteSize)
		n, err := r.Read(barr)
		if err != nil {
			return err
		} else if n != byteSize {
			return fmt.Errorf("bit.Read:short read, expect=%d byte, read=%d byte", byteSize, n)
		}
		err = fillData(barr, order, reflect.Indirect(v), &Offset{})
		if err != io.EOF {
			return err
		}
	default:
		return binary.Read(r, order, data)
	}
	return nil
}
