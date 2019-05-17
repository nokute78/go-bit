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

// Package bit provides functions for bit.
package bit

import (
	"fmt"
)

// Offset represents offset to access bits in byte slices.
type Offset struct {
	Byte uint64 /* Offset in byte. */
	Bit  uint64 /* Offset in bit.  */
}

// Normalize updates off.Byte if off.Bit >= 8.
// e.g. Offset{Byte: 3, Bit: 53} -> Offset{Byte: 9, Bit: 5}
func (off *Offset) Normalize() {
	if off.Bit < 8 {
		return
	}
	b := off.Bit / 8
	off.Byte = off.Byte + b
	off.Bit = off.Bit % 8
}

// GetBit returns 1 or 0.
// GetBit reads b at Offset off, returns the bit.
func GetBit(b []byte, off Offset) (byte, error) {
	a, err := GetBitNotShift(b, off)
	if a > 0x0 {
		return 0x1, err
	}
	return 0x0, err
}

// GetBitNotShift reads b at Offset off, returns the bit.
// Return value is not bit shifted.
func GetBitNotShift(b []byte, off Offset) (byte, error) {
	off.Normalize()
	if len(b) <= int(off.Byte) {
		return 0x0, fmt.Errorf("out of range")
	}
	return b[off.Byte] & (1 << off.Bit), nil
}

// Compare returns an integer comparing two Offsets.
// The result will be 0 if off==b, -1 if off < b, and +1 if off > b.
func (off Offset) Compare(b Offset) int {
	off.Normalize()
	b.Normalize()

	if off.Byte > b.Byte {
		return 1
	} else if off.Byte < b.Byte {
		return -1
	}

	/* off.Byte == b.Byte */
	if off.Bit > b.Bit {
		return 1
	} else if off.Bit < b.Bit {
		return -1
	}
	return 0
}

// OffsetInBit returns offset in bit.
// e.g. Offset{Byte:3, Bit:2} -> 26.
func (off Offset) OffsetInBit() uint64 {
	return off.Byte*8 + off.Bit
}

// AddOffset adds diff and returns new Offset.
func (off Offset) AddOffset(diff Offset) (Offset, error) {
	ret := Offset{Byte: off.Byte + diff.Byte, Bit: off.Bit + diff.Bit}
	ret.Normalize()

	return ret, nil
}

// SubOffset subs diff and returns new Offset.
// diff must be larger then off.
func (off Offset) SubOffset(diff Offset) (Offset, error) {
	if off.Compare(diff) < 0 {
		return Offset{}, fmt.Errorf("negative")
	}

	ret := Offset{Byte: 0, Bit: off.OffsetInBit() - diff.OffsetInBit()}
	ret.Normalize()

	return ret, nil
}

// GetBits returns byte slice.
// GetBits reads bytes slice from Offset off. Read size is bitSize in bit.
func GetBits(bytes []byte, off Offset, bitSize uint64) (ret []byte, err error) {
	tail, err := off.AddOffset(Offset{Byte: 0, Bit: bitSize})
	if err != nil {
		return []byte{}, err
	}
	if len(bytes) <= int(tail.Byte) {
		return []byte{}, fmt.Errorf("out of range")
	}

	var retSize uint64
	if bitSize%8 > 0 {
		retSize = bitSize/8 + 1
	} else {
		retSize = bitSize
	}
	ret = make([]byte, retSize)

	for i := uint64(0); i < bitSize; i++ {
		bitOff, err := off.AddOffset(Offset{0, i})
		if err != nil {
			return []byte{}, err
		}
		bit, err := GetBit(bytes, bitOff)
		if err != nil {
			return []byte{}, err
		}
		ret[i/8] = ret[i/8] | (bit << (i % 8))
	}
	return ret, nil
}
