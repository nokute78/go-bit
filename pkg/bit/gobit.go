    
/*
   Copyright 2019 Takahiro Yamashita
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
	"fmt"
)

type Offset struct {
	Byte uint64 /* Offset in byte. */
	Bit  uint64 /* Offset in bit.  */
}

func (off *Offset) Normalize() {
	b := off.Bit / 8
	off.Byte = off.Byte + b
	off.Bit = off.Bit % 8
}

func GetBit(b []byte, off Offset) (byte, error) {
	a, err := GetBitNotShift(b, off)
	if a > 0x0 {
		return 0x1, err
	} else {
		return 0x0, err
	}

}

func GetBitNotShift(b []byte, off Offset) (byte, error) {
	off.Normalize()
	if len(b) <= int(off.Byte) {
		return 0x0, fmt.Errorf("out of range")
	}
	return b[off.Byte] & (1 << off.Bit), nil
}

func Compare(a, b Offset) int {
	a.Normalize()
	b.Normalize()

	if a.Byte > b.Byte {
		return 1
	} else if a.Byte < b.Byte {
		return -1
	}

	/* a.Byte == b.Byte */
	if a.Bit > b.Bit {
		return 1
	} else if a.Bit < b.Bit {
		return -1
	}
	return 0
}

func (off Offset) SizeInBit() uint64 {
	return off.Byte*8 + off.Bit
}

func (off Offset) AddOffset(diff Offset) (Offset, error) {
	ret := Offset{Byte: off.Byte + diff.Byte, Bit: off.Bit + diff.Bit}
	ret.Normalize()

	return ret, nil
}

func (off Offset) SubOffset(diff Offset) (Offset, error) {
	if Compare(off, diff) < 0 {
		return Offset{}, fmt.Errorf("negative")
	}

	ret := Offset{Byte: 0, Bit: off.SizeInBit() - diff.SizeInBit()}
	ret.Normalize()

	return ret, nil
}

func GetBits(bytes []byte, off Offset, bit_size uint64) (ret []byte, err error) {
	tail, err := off.AddOffset(Offset{Byte: 0, Bit: bit_size})
	if err != nil {
		return []byte{}, err
	}
	if len(bytes) <= int(tail.Byte) {
		return []byte{}, fmt.Errorf("out of range")
	}

	var ret_size uint64
	if bit_size%8 > 0 {
		ret_size = bit_size/8 + 1
	} else {
		ret_size = bit_size
	}
	ret = make([]byte, ret_size)

	for i := uint64(0); i < bit_size; i++ {
		bit_off, err := off.AddOffset(Offset{0, i})
		if err != nil {
			return []byte{}, err
		}
		bit, err := GetBit(bytes, bit_off)
		if err != nil {
			return []byte{}, err
		}
		ret[i/8] = ret[i/8] | (bit << (i % 8))
	}
	return ret, nil
}
