    
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
	"bytes"
	"testing"
)

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

func TestGetBit(t *testing.T) {
	type testcase struct {
		name     string
		bytes    []byte
		off      Offset
		expected byte
	}

	cases := []testcase{
		{"0000_1000[2]", []byte{0x08}, Offset{0, 2}, 0},
		{"0000_1000[3]", []byte{0x08}, Offset{0, 3}, 1},
		{"0000_0100_0000_0000[8+2]", []byte{0x00, 0x04}, Offset{1, 2}, 1},
		{"0000_0100_0000_0000[8+1]", []byte{0x00, 0x04}, Offset{1, 1}, 0},
	}

	for _, v := range cases {
		b, err := GetBit(v.bytes, v.off)
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
		_, err := GetBit(v.bytes, v.off)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}

func TestGetBitNotShift(t *testing.T) {
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
		b, err := GetBitNotShift(v.bytes, v.off)
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
		_, err := GetBitNotShift(v.bytes, v.off)
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
		ret := Compare(v.a, v.b)
		if ret != v.expected {
			t.Errorf("%s: mismatch. given %d. expected %d", v.name, ret, v.expected)
		}
	}
}

func TestSizeInBit(t *testing.T) {
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
		ret := v.a.SizeInBit()
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

func TestGetBits(t *testing.T) {
	type testcase struct {
		name string
		Offset
		bit_size uint64
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
		ret, err := GetBits(v.testdata, v.Offset, v.bit_size)
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
		_, err := GetBits(v.testdata, v.Offset, v.bit_size)
		if err == nil {
			t.Errorf("%s: It should be error", v.name)
		}
	}
}
