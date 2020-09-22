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
	"bytes"
	"encoding/binary"
	"github.com/nokute78/go-bit/pkg/bit/v2"
	"testing"
)

func TestWritePrimitive(t *testing.T) {
	// uint8
	buf := bytes.NewBuffer([]byte{})
	var u8 uint8 = 0xbb
	if err := bit.Write(buf, binary.LittleEndian, u8); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret, err := buf.ReadByte()
	if err != nil {
		t.Errorf("ReadByte err=%s", err)
	} else if ret != u8 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret, u8)
	}

	// uint16
	buf.Reset()
	var u16 uint16 = 0xbbee
	if err := bit.Write(buf, binary.LittleEndian, u16); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret16 := binary.LittleEndian.Uint16(buf.Bytes())
	if ret16 != u16 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret16, u16)
	}

	// uint32
	buf.Reset()
	var u32 uint32 = 0xbbeeccff
	if err := bit.Write(buf, binary.LittleEndian, u32); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret32 := binary.LittleEndian.Uint32(buf.Bytes())
	if ret32 != u32 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret32, u32)
	}
}

func TestWriteByteSlice(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	input := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	if err := bit.Write(buf, binary.LittleEndian, input); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret := buf.Bytes()
	if bytes.Compare(ret, input) != 0 {
		t.Errorf("mismatch\n given:%x\n expect:%x", ret, input)
	}

	// bigendian
	buf.Reset()
	if err := bit.Write(buf, binary.BigEndian, input); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret = buf.Bytes()
	if len(ret) != len(input) {
		t.Errorf("len error given=%d expect=%d", len(ret), len(input))
	}

	length := len(ret)
	expect := make([]byte, length)
	for i := 0; i < length/2; i++ {
		expect[i], expect[length-i-1] = input[length-i-1], input[i]
	}
	if bytes.Compare(expect, ret) != 0 {
		t.Errorf("mismatch\n given=%x expect=%x", ret, expect)
	}

}

func TestWriteByteArray(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	input := [4]byte{0xaa, 0xbb, 0xcc, 0xdd}
	if err := bit.Write(buf, binary.LittleEndian, input); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret := buf.Bytes()
	iexpect := input[:]
	if bytes.Compare(ret, iexpect) != 0 {
		t.Errorf("mismatch\n given:%x\n expect:%x", ret, iexpect)
	}

	// bigendian
	buf.Reset()
	if err := bit.Write(buf, binary.BigEndian, input); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret = buf.Bytes()
	if len(ret) != len(input) {
		t.Errorf("len error given=%d expect=%d", len(ret), len(input))
	}

	length := len(ret)
	expect := make([]byte, length)
	for i := 0; i < length/2; i++ {
		expect[i], expect[length-i-1] = input[length-i-1], input[i]
	}
	if bytes.Compare(expect, ret) != 0 {
		t.Errorf("mismatch\n given=%x expect=%x", ret, expect)
	}

}
