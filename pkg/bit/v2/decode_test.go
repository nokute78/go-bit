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
	"io"
	"testing"
)

func TestDataSizeInBits(t *testing.T) {
	type testcase struct {
		name   string
		input  interface{}
		expect int
	}

	type S72 struct {
		B byte
		I int64
	}
	type Nest struct {
		B byte
		S S72
	}

	cases := []testcase{
		{"Bin Slice", bit.NewBits(33, true), 33},
		{"int", int64(100), 64},
		{"struct", S72{}, 72},
		{"slice", []byte{1, 2}, 16},
		{"array struct", []S72{{1, 1}, {2, 2}}, 144},
		{"nest struct", Nest{}, 80},
	}

	for _, v := range cases {
		i := bit.Size(v.input)
		if i != v.expect {
			t.Errorf("%s:given=%d expect=%d", v.name, i, v.expect)
		}
	}
}

func TestReadPrimitive(t *testing.T) {
	// uint8
	br := bytes.NewReader([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
	var u8 uint8
	if err := bit.Read(br, binary.LittleEndian, &u8); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err := br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u8 != 0 {
		t.Errorf("uint8:given %x expect 0", u8)
	}

	// uint16
	var u16 uint16
	if err := bit.Read(br, binary.LittleEndian, &u16); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err = br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u16 != 0x100 {
		t.Errorf("uint16:given %x expect 0", u16)
	}

	// uint32
	var u32 uint32
	if err := bit.Read(br, binary.LittleEndian, &u32); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err = br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u32 != 0x03020100 {
		t.Errorf("uint32:given %x expect 0", u32)
	}

	// uint64
	var u64 uint64
	if err := bit.Read(br, binary.LittleEndian, &u64); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err = br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u64 != 0x0706050403020100 {
		t.Errorf("uint64:given %x expect 0", u64)
	}
}

func TestReadStruct(t *testing.T) {
	type Sample struct {
		Header   byte
		Reserved [16]bit.Bit
		Id       [4]bit.Bit
		Rev      [4]bit.Bit
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewReader([]byte{0x7f, 0xff, 0xff, 0x51, 0xaa, 0xbb, 0xcc, 0xdd})
	if err := bit.Read(br, binary.LittleEndian, &s); err != nil {
		t.Fatalf("error:%s", err)
	}

	if s.Header != 0x7f {
		t.Errorf("header: given=%x expect=%x", s.Header, 0x7f)
	}
	for i := 0; i < len(s.Reserved); i++ {
		if !s.Reserved[i] {
			t.Errorf("bit(%d) is not 1", i)
		}
	}
	b := bit.BitsToBytes(s.Id[:], binary.LittleEndian)
	if len(b) != 1 {
		t.Errorf("Id size error")
	}
	if b[0] != 0x1 {
		t.Errorf("Id: given=%x expect=%x", b[0], 0x1)
	}

	b = bit.BitsToBytes(s.Rev[:], binary.LittleEndian)
	if len(b) != 1 {
		t.Errorf("Rev size error")
	}
	if b[0] != 0x5 {
		t.Errorf("Rev: given=%x expect=%x", b[0], 0x5)
	}

	expect := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	if len(s.Data) != len(expect) {
		t.Errorf("Data size error")
	}
	if bytes.Compare(expect, s.Data[:]) != 0 {
		t.Errorf("Data: given=%v expect=%v", s.Data, expect)
	}
}

func TestStructTag(t *testing.T) {
	type Sample struct {
		Bit      [4]bit.Bit
		Reserved [4]bit.Bit `bit:"-"` // ignored
		Val      byte
	}

	s := Sample{}
	br := bytes.NewReader([]byte{0xff, 0xaa})
	if err := bit.Read(br, binary.LittleEndian, &s); err != nil {
		t.Fatalf("error:%s\n", err)
	}

	// ignored field
	for i, v := range s.Reserved {
		if v {
			t.Errorf("%d:bit is 1!?", i)
		}
	}

	for i, v := range s.Bit {
		if !v {
			t.Errorf("%d:bit is 0!?", i)
		}
	}
	if s.Val != 0xaf {
		t.Errorf("given=0x%x expect=0xaf", s.Val)
	}

	// skip case
	type Sample2 struct {
		Bit      [4]bit.Bit
		Reserved [4]bit.Bit `bit:"skip"` // skip
		Val      byte
	}

	s2 := Sample2{}
	br = bytes.NewReader([]byte{0xff, 0xaa})
	if err := bit.Read(br, binary.LittleEndian, &s2); err != nil {
		t.Fatalf("error:%s\n", err)
	}

	// ignored field
	for i, v := range s2.Reserved {
		if v {
			t.Errorf("%d:bit is 1!?", i)
		}
	}

	for i, v := range s2.Bit {
		if !v {
			t.Errorf("%d:bit is 0!?", i)
		}
	}
	if s2.Val != 0xaa {
		t.Errorf("given=0x%x expect=0xaa", s2.Val)
	}
}

func TestReadBigEndian(t *testing.T) {
	type TcpHeader struct {
		SrcPort    uint16
		DstPort    uint16
		SeqNo      uint32
		AckNo      uint32
		HeaderLen  [4]bit.Bit
		Reserved   [3]bit.Bit
		NS         bit.Bit
		CWR        bit.Bit
		ECE        bit.Bit
		URG        bit.Bit
		ACK        bit.Bit
		PSH        bit.Bit
		RST        bit.Bit
		SYN        bit.Bit
		FIN        bit.Bit
		WinSize    uint16
		CheckSum   uint16
		EmePointer uint16
	}

	s := TcpHeader{}
	br := bytes.NewReader([]byte{0xd8, 0x65, 0x01, 0xbb, 0x4b, 0xe0, 0x76, 0xcd, 0x48, 0xc8, 0x70, 0x8f, 0x50, 0x10, 0x10,
		0x18, 0x0e, 0xc1, 0x00, 0x00})
	if err := bit.Read(br, binary.BigEndian, &s); err != nil {
		t.Fatalf("error:%s", err)
	}
	if s.SrcPort != 0xd865 {
		t.Errorf("SrcPort:given=0x%x expect=0x%x", s.SrcPort, 0xd865)
	}
	if s.DstPort != 0x1bb {
		t.Errorf("DstPort:given=0x%x expect=0x%x", s.DstPort, 0x1bb)
	}
	if !s.ACK {
		t.Errorf("ACK is false")
	}
	if s.HeaderLen[0] || !s.HeaderLen[1] || s.HeaderLen[2] || !s.HeaderLen[3] {
		t.Errorf("HeaderLength is not 5. %v\n", s.HeaderLen)
	}
}

func BenchmarkReadStruct(b *testing.B) {
	type Sample struct {
		Header   byte
		Reserved [16]bit.Bit
		Id       [4]bit.Bit
		Rev      [4]bit.Bit
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewReader([]byte{0x7f, 0xff, 0xff, 0x51, 0xaa, 0xbb, 0xcc, 0xdd})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := br.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
		err = bit.Read(br, binary.LittleEndian, &s)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
	}
}
