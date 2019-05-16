package gobit

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

func GetBits(bytes []byte, off *Offset, bit_size uint64) (ret []byte, err error) {
	if off == nil {
		return []byte{}, fmt.Errorf("offset is nil")
	}

	tail, err := off.AddOffset(Offset{Byte: 0, Bit: bit_size})
	if err != nil {
		return []byte{}, err
	}

	tail.Normalize()
	if len(bytes) <= int(tail.Byte) {
		return []byte{}, fmt.Errorf("out of range")
	}

	length, err := tail.SubOffset(*off)
	if err != nil {
		return []byte{}, err
	}

	if length.Bit > 0 {
		ret = make([]byte, length.Byte+1)
	} else {
		ret = make([]byte, length.Byte)
	}

	/*
		partBytes := bytes[off.Byte: tail.Byte]

		for i, v := range partBytes {
			ret[i]
		}
	*/

	return []byte{}, nil
}
