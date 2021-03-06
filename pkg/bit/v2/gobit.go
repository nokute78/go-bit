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

// Package bit provides functions for bit.
package bit

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrOutOfRange = errors.New("out of range")
)

type Bit bool

func (b Bit) String() string {
	if b {
		return "1"
	}
	return "0"
}

// BitsToBytes converts the unit. bit -> byte.
func BitsToBytes(b []Bit, o binary.ByteOrder) []byte {
	size := SizeInByte(b)
	ret := make([]byte, size)

	bitc := 0
	idx := 0
	if o == binary.BigEndian {
		idx = len(ret) - 1
	}
	for i := 0; i < len(b); i++ {
		if b[i] {
			ret[idx] = ret[idx] | (1 << bitc)
		}
		bitc += 1
		if bitc == 8 {
			bitc = 0
			if o == binary.BigEndian {
				idx -= 1
			} else {
				idx += 1
			}
		}
	}

	return ret
}

// GetBitsBitEndian returns Bit slice.
// If order is LittleEndian, it is same as GetBits function.
// It respect bit order endianness when order is BigEndian.
func GetBitsBitEndian(b []byte, o Offset, bitSize uint64, order binary.ByteOrder) ([]Bit, error) {
	_, err := isInRange(b, o, bitSize)
	if err != nil {
		return []Bit{}, err
	}

	if order == binary.BigEndian {
		return getBitsBigBitEndian(b, o, bitSize)
	}
	return GetBits(b, o, bitSize, order)
}

// this function treats bit order endian.
//  e.g. b = []byte{0x50} and bitSize is 4. It returns []Bit{true, false, true, false} = 0x5
//    0x50  =   0101|0000
//              3210|----   ( Offset when bitSize=4 )
//              7654|3210   ( Offset when bitSize=8 )
func getBitsBigBitEndian(b []byte, o Offset, bitSize uint64) ([]Bit, error) {
	ret := make([]Bit, bitSize)
	byteAddr := o.Byte
	bitAddr := 7 - int(o.Bit)

	for i := 0; i < int(bitSize); i++ {
		if b[byteAddr]&(1<<bitAddr) > 0 {
			ret[len(ret)-1-i] = true
		} else {
			ret[len(ret)-1-i] = false
		}
		bitAddr -= 1
		if bitAddr == -1 {
			byteAddr += 1
			bitAddr = 7
		}
	}
	return ret, nil
}

// BytesToBits returns Bit slices. bitSize is the size of Bit slice.
func BytesToBits(b []byte, bitSize uint64, o binary.ByteOrder) ([]Bit, error) {
	if o == binary.BigEndian {
		return getBitsBigBitEndian(b, Offset{}, bitSize)
	}
	return GetBits(b, Offset{}, bitSize, binary.LittleEndian)
}

// SizeInByte returns size of []Bit slice in byte.
// e.g. It returns  len([]Bit) == 9.
func SizeInByte(b []Bit) int {
	return sizeOfBits(len(b))
}

func sizeOfBits(bitsize int) int {
	ret := bitsize / 8
	if bitsize%8 > 0 {
		ret += 1
	}
	return ret
}

// Offset represents offset to access bits in byte slices.
type Offset struct {
	Byte uint64 /* Offset in byte. */
	Bit  uint64 /* Offset in bit.  */
}

func (o Offset) String() string {
	return fmt.Sprintf("[Byte:%d,Bit:%d]", o.Byte, o.Bit)
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

// Bits returns offset in bit.
// e.g. Offset{Byte:3, Bit:2} -> 26.
func (off Offset) Bits() uint64 {
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

	ret := Offset{Byte: 0, Bit: off.Bits() - diff.Bits()}
	ret.Normalize()

	return ret, nil
}

func checkRange(b []byte, off Offset) error {
	if len(b) < int(off.Byte) || (len(b) == int(off.Byte) && int(off.Bit) > 0) {
		return ErrOutOfRange
	}
	return nil
}

// SetBit sets bit on b at off.
// Bit is 0 if val == 0, 1 if val > 0.
// SetBit returns error if error occurred.
func SetBit(b []byte, off Offset, val Bit, o binary.ByteOrder) error {
	off.Normalize()
	if err := checkRange(b, off); err != nil {
		return fmt.Errorf("SetBit:%w", err)
	}

	addr := int(off.Byte)
	if o == binary.BigEndian {
		addr = len(b) - 1 - addr
	}
	if val {
		/* set bit */
		b[addr] |= 1 << off.Bit
	} else {
		/* clear bit */
		b[addr] &= ^(1 << off.Bit)
	}
	return nil
}

// GetBit returns 1 or 0.
// GetBit reads b at Offset off, returns the bit.
func GetBit(b []byte, off Offset, o binary.ByteOrder) (Bit, error) {
	v, err := GetBitAsByte(b, off, o)
	if err != nil {
		return false, err
	}
	if v > 0 {
		return true, nil
	}
	return false, nil
}

// GetBitAsByte returns byte (1 or 0).
// GetBitAsByte reads b at Offset off, returns the bit.
func GetBitAsByte(b []byte, off Offset, o binary.ByteOrder) (byte, error) {
	a, err := GetBitAsByteNotShift(b, off, o)
	if a > 0x0 {
		return 0x1, err
	}
	return 0x0, err
}

// GetBitAsByteNotShift reads b at Offset off, returns the bit.
// Return value is not bit shifted.
func GetBitAsByteNotShift(b []byte, off Offset, o binary.ByteOrder) (byte, error) {
	off.Normalize()
	byteAddr := int(off.Byte)
	bitAddr := int(off.Bit)
	if o == binary.BigEndian {
		byteAddr = len(b) - 1 - byteAddr
		//		bitAddr = 7 - bitAddr
	}
	if err := checkRange(b, off); err != nil {
		return 0x0, fmt.Errorf("GetBitAsByteNotShift:%w", err)
	}
	return b[byteAddr] & (1 << bitAddr), nil
}

func isInRange(b []byte, off Offset, bitSize uint64) (bool, error) {
	tail, err := off.AddOffset(Offset{Byte: 0, Bit: bitSize})
	if err != nil {
		return false, err
	}
	if err := checkRange(b, tail); err != nil {
		return false, fmt.Errorf("isInRange:%w", err)
	}
	return true, nil
}

// SetBits sets bits on bytes at off.
// The length to set is bitSize.
// SetBits returns error if error occurred.
func SetBits(bytes []byte, off Offset, setBits []Bit, o binary.ByteOrder) error {
	bitSize := uint64(len(setBits))
	sb := BitsToBytes(setBits, o)

	_, err := isInRange(bytes, off, bitSize)
	if err != nil {
		return err
	}

	for i := uint64(0); i < bitSize; i++ {
		bit, err := GetBit(sb, Offset{0, i}, o)
		if err != nil {
			return err
		}
		err = SetBit(bytes, Offset{off.Byte, off.Bit + i}, bit, o)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetBitsBitEndian sets bits in b.
// If order is LittleEndian, it is same as GetBits function.
// It respect bit order endianness when order is BigEndian.
func SetBitsBitEndian(b []byte, off Offset, setBits []Bit, order binary.ByteOrder) error {
	bitSize := uint64(len(setBits))
	_, err := isInRange(b, off, bitSize)
	if err != nil {
		return err
	}

	if order == binary.BigEndian {
		return setBitsBigBitEndian(b, off, setBits, bitSize)
	}
	return SetBits(b, off, setBits, order)
}

// this function treats bit order endian.
//    0x50  =   0101|0000
//              3210|----   ( Offset when bitSize=4 )
//              7654|3210   ( Offset when bitSize=8 )
func setBitsBigBitEndian(b []byte, off Offset, setBits []Bit, bitSize uint64) error {
	byteAddr := off.Byte
	bitAddr := 7 - int(off.Bit)

	for i := 0; i < int(bitSize); i++ {
		if setBits[int(bitSize)-1-i] {
			b[byteAddr] |= 1 << bitAddr
		} else {
			b[byteAddr] &= ^(1 << bitAddr)
		}

		bitAddr -= 1
		if bitAddr == -1 {
			byteAddr += 1
			bitAddr = 7
		}
	}

	return nil
}

// GetBits returns Bit slice.
// GetBits reads bytes slice from Offset off. Read size is bitSize in bit.
func GetBits(bytes []byte, off Offset, bitSize uint64, o binary.ByteOrder) (ret []Bit, err error) {
	ok, err := isInRange(bytes, off, bitSize)
	if !ok || err != nil {
		return []Bit{}, err
	}
	ret = make([]Bit, bitSize)

	for i := 0; i < int(bitSize); i++ {
		bitOff, err := off.AddOffset(Offset{0, uint64(i)})
		if err != nil {
			return []Bit{}, err
		}
		bit, err := GetBit(bytes, bitOff, o)
		if err != nil {
			return []Bit{}, err
		}
		ret[i] = bit

	}
	return ret, nil
}

// GetBitsAsByte returns byte slice.
// GetBitsAsByte reads bytes slice from Offset off. Read size is bitSize in bit.
func GetBitsAsByte(bytes []byte, off Offset, bitSize uint64, o binary.ByteOrder) (ret []byte, err error) {
	ok, err := isInRange(bytes, off, bitSize)
	if !ok || err != nil {
		return []byte{}, err
	}
	bits := make([]Bit, bitSize)

	for i := uint64(0); i < bitSize; i++ {
		bitOff, err := off.AddOffset(Offset{0, i})
		if err != nil {
			return []byte{}, err
		}
		bit, err := GetBit(bytes, bitOff, o)
		if err != nil {
			return []byte{}, err
		}
		if bit {
			bits[i] = true
		}
	}

	return BitsToBytes(bits, o), nil
}

// NewBits generates slice of Bit.
func NewBits(size uint64, v Bit) []Bit {
	b := make([]Bit, size)
	if v {
		for i := 0; i < int(size); i++ {
			b[i] = true
		}
	}
	return b
}
