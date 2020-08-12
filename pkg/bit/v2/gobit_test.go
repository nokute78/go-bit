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
	"bytes"
	"encoding/binary"
	"testing"
)

func TestBitsToBytes(t *testing.T) {
	type testcase struct {
		name   string
		input  []Bit
		expect []byte
	}

	cases := []testcase{
		{"normal", []Bit{false, true, true, true}, []byte{0xe}},
		{"2byte", []Bit{false, true, true, true, false, false, false, false, true}, []byte{0xe, 0x1}},
	}

	for _, v := range cases {
		ret := BitsToBytes(v.input, binary.LittleEndian)
		if bytes.Compare(ret, v.expect) != 0 {
			t.Errorf("%s: given=%v expect=%v", v.name, ret, v.expect)
		}
	}
}

func TestSizeOfBits(t *testing.T) {
	type testcase struct {
		name   string
		input  []Bit
		expect int
	}

	cases := []testcase{
		{"8bit", NewBits(8, true), 1},
		{"9bit", NewBits(9, true), 2},
		{"0bit", []Bit{}, 0},
	}
	for _, v := range cases {
		ret := SizeOfBits(v.input)
		if ret != v.expect {
			t.Errorf("%s:given=%d expect=%d", v.name, ret, v.expect)
		}
	}
}

func TestOffsetNormalize(t *testing.T) {
	type testcase struct {
		name   string
		before Offset
		after  Offset
	}

	cases := []testcase{
		{"2byte0bit", Offset{Byte: 2, Bit: 0}, Offset{Byte: 2, Bit: 0}},
		{"2byte1bit", Offset{Byte: 2, Bit: 1}, Offset{Byte: 2, Bit: 1}},
		{"0byte16bit", Offset{Byte: 0, Bit: 16}, Offset{Byte: 2, Bit: 0}},
		{"1byte9bit", Offset{Byte: 1, Bit: 9}, Offset{Byte: 2, Bit: 1}},
	}

	for _, v := range cases {
		v.before.Normalize()
		if v.before != v.after {
			t.Errorf("%s: mismatch. given %v. expected %v", v.name, v.before, v.after)
		}
	}
}

func TestSetBit(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		bit      Bit
		expected []byte
	}

	cases := []testcase{
		{"0000_0000[2] on", []byte{0x00}, Offset{0, 2}, true, []byte{0x04}},
		{"0000_0000[2] off", []byte{0x00}, Offset{0, 2}, false, []byte{0x00}},
		{"1111_1111[2] on", []byte{0xff}, Offset{0, 2}, true, []byte{0xff}},
		{"1111_1111[2] off", []byte{0xff}, Offset{0, 2}, false, []byte{0xfb}},
		{"0000_0000_0000_0000[9] on", []byte{0x00, 0x00}, Offset{0, 9}, true, []byte{0x00, 0x02}},
		{"0000_0000_0000_0000[9] off", []byte{0x00, 0x00}, Offset{0, 9}, false, []byte{0x00, 0x00}},
		{"1111_1111_1111_1111[9] on", []byte{0xff, 0xff}, Offset{0, 9}, true, []byte{0xff, 0xff}},
		{"1111_1111_1111_1111[9] off", []byte{0xff, 0xff}, Offset{0, 9}, false, []byte{0xff, 0xfd}},
	}

	for _, v := range cases {
		err := SetBit(v.bytes, v.off, v.bit, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(v.bytes, v.expected) != 0 {
			t.Errorf("%s: mismatch. given 0x%x. expected 0x%x", v.name, v.bytes, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x0}, Offset{128, 0}, false, []byte{0x0}},
	}

	for _, v := range errcases {
		err := SetBit(v.bytes, v.off, v.bit, binary.LittleEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}

}

func TestSetBitBigEndian(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		bit      Bit
		expected []byte
	}

	cases := []testcase{
		{"0000_0000[2] on", []byte{0x00}, Offset{0, 2}, true, []byte{0x04}},
		{"0000_0000[2] off", []byte{0x00}, Offset{0, 2}, false, []byte{0x00}},
		{"1111_1111[2] on", []byte{0xff}, Offset{0, 2}, true, []byte{0xff}},
		{"1111_1111[2] off", []byte{0xff}, Offset{0, 2}, false, []byte{0xfb}},
		{"0000_0000_0000_0000[9] on", []byte{0x00, 0x00}, Offset{0, 9}, true, []byte{0x02, 0x00}},
		{"0000_0000_0000_0000[9] off", []byte{0x00, 0x00}, Offset{0, 9}, false, []byte{0x00, 0x00}},
		{"1111_1111_1111_1111[9] on", []byte{0xff, 0xff}, Offset{0, 9}, true, []byte{0xff, 0xff}},
		{"1111_1111_1111_1111[9] off", []byte{0xff, 0xff}, Offset{0, 9}, false, []byte{0xfd, 0xff}},
	}

	for _, v := range cases {
		err := SetBit(v.bytes, v.off, v.bit, binary.BigEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(v.bytes, v.expected) != 0 {
			t.Errorf("%s: mismatch. given 0x%x. expected 0x%x", v.name, v.bytes, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x0}, Offset{128, 0}, false, []byte{0x0}},
	}

	for _, v := range errcases {
		err := SetBit(v.bytes, v.off, v.bit, binary.BigEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}

}

func TestGetBit(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		expected Bit
	}

	cases := []testcase{
		{"0000_1000[2]", []byte{0x08}, Offset{0, 2}, false},
		{"0000_1000[3]", []byte{0x08}, Offset{0, 3}, true},
		{"0000_0100_0000_0000[8+2]", []byte{0x00, 0x04}, Offset{1, 2}, true},
		{"0000_0100_0000_0000[8+1]", []byte{0x00, 0x04}, Offset{1, 1}, false},
	}

	for _, v := range cases {
		b, err := GetBit(v.bytes, v.off, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if b != v.expected {
			t.Errorf("%s: mismatch. given %t. expected %t", v.name, b, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x0}, Offset{128, 0}, false},
	}

	for _, v := range errcases {
		_, err := GetBit(v.bytes, v.off, binary.LittleEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitBigEndian(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		expected Bit
	}

	cases := []testcase{
		{"0000_1000[3]", []byte{0x08}, Offset{0, 3}, true},
		{"0000_1000[4]", []byte{0x08}, Offset{0, 4}, false},
		{"0000_0000_0000_0001[0]", []byte{0x00, 0x01}, Offset{0, 0}, true},
		{"0000_0000_0000_0001[1]", []byte{0x00, 0x01}, Offset{0, 1}, false},
	}

	for _, v := range cases {
		b, err := GetBit(v.bytes, v.off, binary.BigEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if b != v.expected {
			t.Errorf("%s: mismatch. given %t. expected %t", v.name, b, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x0}, Offset{128, 0}, false},
	}

	for _, v := range errcases {
		_, err := GetBit(v.bytes, v.off, binary.BigEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitAsByteNotShift(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		expected byte
	}

	cases := []testcase{
		{"0000_1000[2]", []byte{0x08}, Offset{0, 2}, 0x0},
		{"0000_1000[3]", []byte{0x08}, Offset{0, 3}, 0x08},
		{"0000_0100_0000_0000[8+2]", []byte{0x00, 0x04}, Offset{1, 2}, 0x04},
		{"0000_0100_0000_0000[8+1]", []byte{0x00, 0x04}, Offset{1, 1}, 0x0},
	}

	for _, v := range cases {
		b, err := GetBitAsByteNotShift(v.bytes, v.off, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if b != v.expected {
			t.Errorf("%s: mismatch. given 0x%x. expected 0x%x", v.name, b, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x0}, Offset{128, 0}, 0},
	}

	for _, v := range errcases {
		_, err := GetBitAsByteNotShift(v.bytes, v.off, binary.LittleEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitAsByteNotShiftBigEndian(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		expected byte
	}

	cases := []testcase{
		{"0000_1000[3]", []byte{0x08}, Offset{0, 3}, 0x8},
		{"0000_1000[4]", []byte{0x08}, Offset{0, 4}, 0x0},
		{"0000_0000_0000_0001[0]", []byte{0x00, 0x01}, Offset{0, 0}, 0x1},
		{"0000_0000_0000_0001[1]", []byte{0x00, 0x01}, Offset{0, 1}, 0x0},
	}

	for _, v := range cases {
		b, err := GetBitAsByteNotShift(v.bytes, v.off, binary.BigEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if b != v.expected {
			t.Errorf("%s: mismatch. given 0x%x. expected 0x%x", v.name, b, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x0}, Offset{128, 0}, 0},
	}

	for _, v := range errcases {
		_, err := GetBitAsByteNotShift(v.bytes, v.off, binary.BigEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestCompare(t *testing.T) {
	type testcase struct {
		name     string
		a        Offset
		b        Offset
		expected int
	}

	cases := []testcase{
		{"a==b", Offset{1, 2}, Offset{1, 2}, 0},
		{"a>b", Offset{1, 2}, Offset{0, 2}, 1},
		{"a<b", Offset{0, 2}, Offset{1, 2}, -1},
		{"a==b(2)", Offset{1, 0}, Offset{0, 8}, 0},
		{"a>b(2)", Offset{1, 10}, Offset{2, 0}, 1},
		{"a<b(2)", Offset{2, 0}, Offset{1, 10}, -1},
	}

	for _, v := range cases {
		ret := v.a.Compare(v.b)
		if ret != v.expected {
			t.Errorf("%s: mismatch. given %d. expected %d", v.name, ret, v.expected)
		}
	}
}

func TestOffsetInBit(t *testing.T) {
	type testcase struct {
		name     string
		a        Offset
		expected uint64
	}

	cases := []testcase{
		{"3bit", Offset{0, 3}, 3},
		{"3byte", Offset{3, 0}, 24},
		{"8+3bit", Offset{1, 3}, 11},
	}

	for _, v := range cases {
		ret := v.a.OffsetInBit()
		if ret != v.expected {
			t.Errorf("%s: mismatch. given %d. expected %d", v.name, ret, v.expected)
		}
	}
}

func TestAddOffset(t *testing.T) {
	type testcase struct {
		name     string
		base     Offset
		diff     Offset
		expected Offset
	}

	cases := []testcase{
		{"bit", Offset{0, 0}, Offset{0, 3}, Offset{0, 3}},
		{"byte", Offset{0, 0}, Offset{3, 0}, Offset{3, 0}},
		{"bit+byte", Offset{0, 0}, Offset{3, 3}, Offset{3, 3}},
		{"normalize", Offset{0, 4}, Offset{2, 7}, Offset{3, 3}},
	}

	for _, v := range cases {
		ret, err := v.base.AddOffset(v.diff)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if ret.Byte != v.expected.Byte {
			t.Errorf("%s: Byte mismatch. given %d. expected %d", v.name, ret.Byte, v.expected.Byte)
		}
		if ret.Bit != v.expected.Bit {
			t.Errorf("%s: Bit mismatch. given %d. expected %d", v.name, ret.Bit, v.expected.Bit)
		}
	}
}

func TestSubOffset(t *testing.T) {
	type testcase struct {
		name     string
		base     Offset
		diff     Offset
		expected Offset
	}

	cases := []testcase{
		{"bit", Offset{0, 5}, Offset{0, 3}, Offset{0, 2}},
		{"byte", Offset{3, 0}, Offset{1, 0}, Offset{2, 0}},
		{"bit+byte", Offset{3, 3}, Offset{1, 2}, Offset{2, 1}},
		{"normalize", Offset{3, 0}, Offset{1, 7}, Offset{1, 1}},
	}

	for _, v := range cases {
		ret, err := v.base.SubOffset(v.diff)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if ret.Byte != v.expected.Byte {
			t.Errorf("%s: Byte mismatch. given %d. expected %d", v.name, ret.Byte, v.expected.Byte)
		}
		if ret.Bit != v.expected.Bit {
			t.Errorf("%s: Bit mismatch. given %d. expected %d", v.name, ret.Bit, v.expected.Bit)
		}
	}
}

func TestGetBitsAsByte(t *testing.T) {
	type testcase struct {
		name string
		Offset
		bitSize  uint64
		testdata []byte
		expected []byte
	}

	cases := []testcase{
		{"from head", Offset{0, 0}, 4, []byte{0x0f}, []byte{0x0f}},
		{"0011_1000", Offset{0, 3}, 3, []byte{0x38}, []byte{0x07}},
		{"0111_1000_0000_0000", Offset{1, 3}, 4, []byte{0x00, 0x78}, []byte{0x0f}},
		{"0000_0011_1100_0000", Offset{0, 6}, 4, []byte{0xc0, 0x03}, []byte{0x0f}},
		{"0111_1111_1100_0000", Offset{0, 6}, 9, []byte{0xc0, 0x7f}, []byte{0xff, 0x01}},
		{"0111_1111_1111_1111_1100_0000", Offset{0, 6}, 17, []byte{0xc0, 0xff, 0x7f}, []byte{0xff, 0xff, 0x01}},
	}

	for _, v := range cases {
		ret, err := GetBitsAsByte(v.testdata, v.Offset, v.bitSize, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(ret, v.expected) != 0 {
			t.Errorf("%s: mismatch. given 0x%x. expected 0x%x", v.name, ret, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", Offset{0, 0}, 128, []byte{0xff}, []byte{}},
	}

	for _, v := range errcases {
		_, err := GetBitsAsByte(v.testdata, v.Offset, v.bitSize, binary.LittleEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitsAsByteBigEndian(t *testing.T) {
	type testcase struct {
		name string
		Offset
		bitSize  uint64
		testdata []byte
		expected []byte
	}

	cases := []testcase{
		{"from head", Offset{0, 0}, 4, []byte{0x0f}, []byte{0x0f}},
		{"0011_1000", Offset{0, 2}, 3, []byte{0x38}, []byte{0x06}},
		{"0111_1000_0000_0000", Offset{0, 11}, 4, []byte{0x78, 0x00}, []byte{0x0f}},
		{"0000_0011_1100_0000", Offset{0, 6}, 4, []byte{0x03, 0xc0}, []byte{0x0f}},
		{"0111_1111_1100_0000", Offset{0, 6}, 9, []byte{0x7f, 0xc0}, []byte{0x01, 0xff}},
		{"0111_1111_1111_1111_1100_0000", Offset{0, 6}, 17, []byte{0x7f, 0xff, 0xc0}, []byte{0x01, 0xff, 0xff}},
	}

	for _, v := range cases {
		ret, err := GetBitsAsByte(v.testdata, v.Offset, v.bitSize, binary.BigEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(ret, v.expected) != 0 {
			t.Errorf("%s: mismatch. given 0x%x. expected 0x%x", v.name, ret, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", Offset{0, 0}, 128, []byte{0xff}, []byte{}},
	}

	for _, v := range errcases {
		_, err := GetBitsAsByte(v.testdata, v.Offset, v.bitSize, binary.BigEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBits(t *testing.T) {
	type testcase struct {
		name string
		Offset
		bitSize  uint64
		testdata []byte
		expected []Bit
	}

	cases := []testcase{
		{"from head", Offset{0, 0}, 4, []byte{0x0f}, NewBits(4, true)},
		{"0011_1000", Offset{0, 3}, 3, []byte{0x38}, NewBits(3, true)},
		{"0111_1000_0000_0000", Offset{1, 3}, 4, []byte{0x00, 0x78}, NewBits(4, true)},
		{"0000_0011_1100_0000", Offset{0, 6}, 4, []byte{0xc0, 0x03}, NewBits(4, true)},
		{"0111_1111_1100_0000", Offset{0, 6}, 9, []byte{0xc0, 0x7f}, NewBits(9, true)},
		{"0111_1111_1111_1111_1100_0000", Offset{0, 6}, 17, []byte{0xc0, 0xff, 0x7f}, NewBits(17, true)},
	}

	for _, v := range cases {
		ret, err := GetBits(v.testdata, v.Offset, v.bitSize, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(BitsToBytes(ret, binary.LittleEndian), BitsToBytes(v.expected, binary.LittleEndian)) != 0 {
			t.Errorf("%s: mismatch. given %t. expected %t", v.name, ret, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", Offset{0, 0}, 128, []byte{0xff}, []Bit{}},
	}

	for _, v := range errcases {
		_, err := GetBits(v.testdata, v.Offset, v.bitSize, binary.LittleEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitsBigEndian(t *testing.T) {
	type testcase struct {
		name string
		Offset
		bitSize  uint64
		testdata []byte
		expected []Bit
	}

	cases := []testcase{
		{"from head", Offset{0, 4}, 4, []byte{0xf0}, NewBits(4, true)},
		{"0011_1000", Offset{0, 3}, 3, []byte{0x38}, NewBits(3, true)},
		{"0111_1000_0000_0000", Offset{0, 11}, 4, []byte{0x78, 0x00}, NewBits(4, true)},
		{"0000_0000_0001_1110", Offset{0, 1}, 4, []byte{0x00, 0x1e}, NewBits(4, true)},
		{"0000_0011_1100_0000", Offset{0, 6}, 4, []byte{0x03, 0xc0}, NewBits(4, true)},
		{"0111_1111_1100_0000", Offset{0, 6}, 9, []byte{0x7f, 0xc0}, NewBits(9, true)},
		{"0111_1111_1111_1111_1100_0000", Offset{0, 6}, 17, []byte{0x7f, 0xff, 0xc0}, NewBits(17, true)},
	}

	for _, v := range cases {
		ret, err := GetBits(v.testdata, v.Offset, v.bitSize, binary.BigEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(BitsToBytes(ret, binary.BigEndian), BitsToBytes(v.expected, binary.BigEndian)) != 0 {
			t.Errorf("%s: mismatch. given %t. expected %t", v.name, ret, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", Offset{0, 0}, 128, []byte{0xff}, []Bit{}},
	}

	for _, v := range errcases {
		_, err := GetBits(v.testdata, v.Offset, v.bitSize, binary.BigEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestSetBits(t *testing.T) {
	type testcase struct {
		name  string
		bytes []byte
		Offset
		bitSize  uint64
		bits     []Bit
		expected []byte
	}

	cases := []testcase{
		{"from head", []byte{0x00}, Offset{0, 0}, 4, NewBits(4, true), []byte{0x0f}},
		{"0011_1000 -> 0000_0000", []byte{0x00}, Offset{0, 3}, 3, NewBits(3, true), []byte{0x38}},
		{"0000_0000 -> 0011_1000", []byte{0x38}, Offset{0, 3}, 3, NewBits(3, false), []byte{0x00}},
		{"0111_1000_0000_0000 -> 0000_0000_0000_0000", []byte{0x00, 0x78}, Offset{1, 3}, 4, NewBits(4, false), []byte{0x00, 0x00}},
		{"0000_0000_0000_0000 -> 0111_1000_0000_0000", []byte{0x00, 0x00}, Offset{1, 3}, 4, NewBits(4, true), []byte{0x00, 0x78}},
		{"0000_0011_1100_0000 -> 0000_0000_0000_0000", []byte{0xc0, 0x03}, Offset{0, 6}, 4, NewBits(4, false), []byte{0x00, 0x00}},
		{"0000_0000_0000_0000 -> 0000_0011_1100_0000", []byte{0x00, 0x00}, Offset{0, 6}, 4, NewBits(4, true), []byte{0xc0, 0x03}},
		{"0111_1111_1100_0000 -> 0000_0000_0000_0000", []byte{0xc0, 0x7f}, Offset{0, 6}, 9, NewBits(9, false), []byte{0x00, 0x00}},
		{"0000_0000_0000_0000 -> 0111_1111_1100_0000", []byte{0x00, 0x00}, Offset{0, 6}, 9, NewBits(9, true), []byte{0xc0, 0x7f}},
		{"0111_1111_1111_1111_1100_0000 -> 0", []byte{0xc0, 0xff, 0x7f}, Offset{0, 6}, 17, NewBits(17, false), []byte{0x00, 0x00, 0x00}},
		{"0 -> 0111_1111_1111_1111_1100_0000", []byte{0x00, 0x00, 0x00}, Offset{0, 6}, 17, NewBits(17, true), []byte{0xc0, 0xff, 0x7f}},
	}

	for _, v := range cases {
		err := SetBits(v.bytes, v.Offset, v.bitSize, v.bits, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(v.bytes, v.expected) != 0 {
			t.Errorf("\n%s: mismatch.\n given 0x%x. expected 0x%x", v.name, v.bytes, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x00}, Offset{0, 0}, 128, NewBits(4, false), []byte{}},
		{"out of range2", []byte{0x00, 0x00}, Offset{0, 0}, 9, NewBits(4, false), []byte{}},
	}

	for _, v := range errcases {
		err := SetBits(v.bytes, v.Offset, v.bitSize, v.bits, binary.LittleEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestSetBitsBigEndian(t *testing.T) {
	type testcase struct {
		name  string
		bytes []byte
		Offset
		bitSize  uint64
		bits     []Bit
		expected []byte
	}

	cases := []testcase{
		{"from head", []byte{0x00}, Offset{0, 4}, 4, NewBits(4, true), []byte{0xf0}},
		{"0000_0000 -> 0011_1000", []byte{0x00}, Offset{0, 3}, 3, NewBits(3, true), []byte{0x38}},
		{"0011_1000 -> 0000_0000", []byte{0x38}, Offset{0, 3}, 3, NewBits(3, false), []byte{0x00}},
		{"0111_1000_0000_0000 -> 0000_0000_0000_0000", []byte{0x78, 0x00}, Offset{0, 11}, 4, NewBits(4, false), []byte{0x00, 0x00}},
		{"0000_0000_0000_0000 -> 0111_1000_0000_0000", []byte{0x00, 0x00}, Offset{0, 11}, 4, NewBits(4, true), []byte{0x78, 0x00}},
		{"0000_0011_1100_0000 -> 0000_0000_0000_0000", []byte{0x03, 0xc0}, Offset{0, 6}, 4, NewBits(4, false), []byte{0x00, 0x00}},
		{"0000_0000_0000_0000 -> 0000_0011_1100_0000", []byte{0x00, 0x00}, Offset{0, 6}, 4, NewBits(4, true), []byte{0x03, 0xc0}},
		{"0111_1111_1100_0000 -> 0000_0000_0000_0000", []byte{0x7f, 0xc0}, Offset{0, 6}, 9, NewBits(9, false), []byte{0x00, 0x00}},
		{"0000_0000_0000_0000 -> 0111_1111_1100_0000", []byte{0x00, 0x00}, Offset{0, 6}, 9, NewBits(9, true), []byte{0x7f, 0xc0}},
		{"0111_1111_1111_1111_1100_0000 -> 0", []byte{0x7f, 0xff, 0xc0}, Offset{0, 6}, 17, NewBits(17, false), []byte{0x00, 0x00, 0x00}},
		{"0 -> 0111_1111_1111_1111_1100_0000", []byte{0x00, 0x00, 0x00}, Offset{0, 6}, 17, NewBits(17, true), []byte{0x7f, 0xff, 0xc0}},
	}

	for _, v := range cases {
		err := SetBits(v.bytes, v.Offset, v.bitSize, v.bits, binary.BigEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if bytes.Compare(v.bytes, v.expected) != 0 {
			t.Errorf("\n%s: mismatch.\n given 0x%x. expected 0x%x", v.name, v.bytes, v.expected)
		}
	}

	errcases := []testcase{
		{"out of range", []byte{0x00}, Offset{0, 0}, 128, NewBits(4, false), []byte{}},
		{"out of range2", []byte{0x00, 0x00}, Offset{0, 0}, 9, NewBits(4, false), []byte{}},
	}

	for _, v := range errcases {
		err := SetBits(v.bytes, v.Offset, v.bitSize, v.bits, binary.BigEndian)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitsAsByteRetSize(t *testing.T) {
	type testcase struct {
		name    string
		off     Offset
		bitsize uint64
		retsize int
	}

	cases := []testcase{
		{"bitsize 7", Offset{Byte: 0, Bit: 0}, 7, 1},
		{"bitsize 8", Offset{Byte: 0, Bit: 0}, 8, 1},
	}

	for _, v := range cases {
		ret, err := GetBitsAsByte([]byte{0xff, 0xff, 0xff, 0xff}, v.off, v.bitsize, binary.LittleEndian)
		if err != nil {
			t.Errorf("%s: Error %s", v.name, err)
		}
		if len(ret) != v.retsize {
			t.Errorf("%s: mismatch. given %d. expected %d", v.name, len(ret), v.retsize)
		}
	}
}

func BenchmarkGetBitsAsByte(b *testing.B) {
	off := Offset{0, 6}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetBitsAsByte([]byte{0xc0, 0xff, 0x7f}, off, 17, binary.LittleEndian)
		if err != nil {
			b.Fatalf("GetBitsAsByte Error!")
		}
	}
}

func BenchmarkSetBits(b *testing.B) {
	off := Offset{0, 6}
	data := NewBits(17, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytes := []byte{0xc0, 0xff, 0x7f}
		if err := SetBits(bytes, off, 17, data, binary.LittleEndian); err != nil {
			b.Fatalf("SetBits Error!")
		}
	}
}
