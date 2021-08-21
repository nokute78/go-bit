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

package bit_test

import (
	"encoding/binary"
	"fmt"
	"github.com/nokute78/go-bit/v2"
)

func ExampleGetBit() {
	b := []byte{0x00, 0x80} /* 1000_0000 0000_0000 in bit */

	off := bit.Offset{Byte: 0, Bit: 15}
	ret, err := bit.GetBit(b, off, binary.LittleEndian)
	if err != nil {
		fmt.Printf("error:%s\n", err)
	}

	fmt.Printf("%t\n", ret)
	// Output:
	// true
}

func ExampleGetBitAsByteNotShift() {
	b := []byte{0x00, 0x80} /* 1000_0000 0000_0000 in bit */

	off := bit.Offset{Byte: 0, Bit: 15}
	ret, err := bit.GetBitAsByteNotShift(b, off, binary.LittleEndian)
	if err != nil {
		fmt.Printf("error:%s\n", err)
	}

	fmt.Printf("0x%x\n", ret)
	// Output:
	// 0x80
}

func ExampleGetBitsAsByte() {
	b := []byte{0x78} /* 0111_1000 in bit */

	/* try to get 4bits(1111b) from 0111_1000 */
	off := bit.Offset{Byte: 0, Bit: 3}

	ret, err := bit.GetBitsAsByte(b, off, 4, binary.LittleEndian)
	if err != nil {
		fmt.Printf("error:%s\n", err)
	}

	fmt.Printf("0x%x\n", ret)
	// Output:
	// 0x0f
}

func ExampleNormalize() {
	off := bit.Offset{Byte: 0, Bit: 17}

	fmt.Printf("Offset: Byte:%d Bit:%d\n", off.Byte, off.Bit)

	off.Normalize()

	fmt.Printf("Offset: Byte:%d Bit:%d\n", off.Byte, off.Bit)
	// Output:
	// Offset: Byte:0 Bit:17
	// Offset: Byte:2 Bit:1
}

func ExampleSetBit() {
	b := []byte{0x00, 0x00} /* 0000_0000 0000_0000 in bit */

	off := bit.Offset{Byte: 0, Bit: 15}
	val := bit.Bit(true)

	err := bit.SetBit(b, off, val, binary.LittleEndian)
	if err != nil {
		fmt.Printf("error:%s\n", err)
	}
	fmt.Printf("0x%x\n", b)
	// Output:
	// 0x0080
}

func ExampleSetBits() {
	b := []byte{0x00, 0x00} /* 0000_0000 0000_0000 in bit */

	off := bit.Offset{Byte: 0, Bit: 8}
	val := []bit.Bit{false, false, false, true} /* 0000_1000 in bit */

	err := bit.SetBits(b, off, val, binary.LittleEndian)
	if err != nil {
		fmt.Printf("error:%s\n", err)
	}

	fmt.Printf("0x%x\n", b)
	// Output:
	// 0x0008
}
