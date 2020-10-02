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

	// uint64
	buf.Reset()
	var u64 uint64 = 0xbbeeccff00112233
	if err := bit.Write(buf, binary.LittleEndian, u64); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret64 := binary.LittleEndian.Uint64(buf.Bytes())
	if ret64 != u64 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret64, u64)
	}
}

func TestWritePrimitiveBE(t *testing.T) {
	// uint8
	buf := bytes.NewBuffer([]byte{})
	var u8 uint8 = 0xbb
	if err := bit.Write(buf, binary.BigEndian, u8); err != nil {
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
	if err := bit.Write(buf, binary.BigEndian, u16); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret16 := binary.BigEndian.Uint16(buf.Bytes())
	if ret16 != u16 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret16, u16)
	}

	// uint32
	buf.Reset()
	var u32 uint32 = 0xbbeeccff
	if err := bit.Write(buf, binary.BigEndian, u32); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret32 := binary.BigEndian.Uint32(buf.Bytes())
	if ret32 != u32 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret32, u32)
	}

	// uint64
	buf.Reset()
	var u64 uint64 = 0xbbeeccff00112233
	if err := bit.Write(buf, binary.BigEndian, u64); err != nil {
		t.Errorf("bit.Write error %s", err)
	}
	ret64 := binary.BigEndian.Uint64(buf.Bytes())
	if ret64 != u64 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret64, u64)
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

func TestWriteStruct(t *testing.T) {
	type S struct {
		B   byte
		U16 uint16
		A   []byte
	}

	s := S{}
	buf := bytes.NewBuffer([]byte{})

	s.B = 0xaa
	s.U16 = 0xbbcc
	s.A = []byte{0xdd, 0xee, 0xff, 0x00, 0x11}

	if err := bit.Write(buf, binary.LittleEndian, s); err != nil {
		t.Errorf("bit.Write err=%s", err)
	}
	expect := []byte{0xaa, 0xcc, 0xbb, 0xdd, 0xee, 0xff, 0x00, 0x11}
	ret := buf.Bytes()
	if bytes.Compare(ret, expect) != 0 {
		t.Errorf("mismatch\n given=%x\n expect=%x", ret, expect)
	}
}

func TestWriteBigEndian(t *testing.T) {
	type TcpHeader struct {
		SrcPort    uint16
		DstPort    uint16
		SeqNo      uint32
		AckNo      uint32
		HeaderLen  [4]bit.Bit /* 12Byte 0-3bit */
		Reserved   [3]bit.Bit /* 12Byte 4-6bit */
		NS         bit.Bit    /* 12Byte 7bit */
		CWR        bit.Bit    /* 13Byte 0bit */
		ECE        bit.Bit    /* 13Byte 1bit */
		URG        bit.Bit    /* 13Byte 2bit */
		ACK        bit.Bit    /* 13Byte 3bit */
		PSH        bit.Bit    /* 13Byte 4bit */
		RST        bit.Bit    /* 13Byte 5bit */
		SYN        bit.Bit    /* 13Byte 6bit */
		FIN        bit.Bit    /* 13Byte 7bit */
		WinSize    uint16
		CheckSum   uint16
		EmePointer uint16
	}

	s := TcpHeader{}
	s.SrcPort = 0xd865
	s.DstPort = 0x1bb
	s.SeqNo = 0x4be076cd
	s.AckNo = 0x48c8708f
	s.ACK = true
	s.HeaderLen = [4]bit.Bit{true, false, true, false}
	s.WinSize = 4120
	s.CheckSum = 0x0ec1

	buf := bytes.NewBuffer([]byte{})
	if err := bit.Write(buf, binary.BigEndian, &s); err != nil {
		t.Fatalf("bit.Write err=%s", err)
	}

	expect := []byte{0xd8, 0x65, 0x01, 0xbb, 0x4b, 0xe0, 0x76, 0xcd, 0x48, 0xc8, 0x70, 0x8f,
		0x50, 0x10, 0x10, 0x18, 0x0e, 0xc1, 0x00, 0x00}
	if bytes.Compare(buf.Bytes(), expect) != 0 {
		t.Errorf("mismatch\n given =%x\n expect=%x", buf.Bytes(), expect)
		t.Errorf("s=%+v", s)
	}
}
