# go-bit
![Go](https://github.com/nokute78/go-bit/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nokute78/go-bit)](https://goreportcard.com/report/github.com/nokute78/go-bit)
[![Go Reference](https://pkg.go.dev/badge/github.com/nokute78/go-bit/v2.svg)](https://pkg.go.dev/github.com/nokute78/go-bit/v2)

A library to read/write bits from a byte slice.

## Installation

```
go get -u github.com/nokute78/go-bit/v2
```

**CAUTION**: `github.com/nokute78/go-bit/pkg/bit/v2` will be deprecated.
Please use `github.com/nokute78/go-bit/v2`.

## Usage

The package supports `binary.Read` like API.

It is an example to decode TCP header. (Big Endian)

1. Save below code as tcp.go.
2. `go mod init go-bit-example`
3. `go mod tidy`
4. `go run tcp.go`

```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/nokute78/go-bit/v2"
)

func main() {
	type TcpHeader struct {
		SrcPort    uint16
		DstPort    uint16
		SeqNo      uint32
		AckNo      uint32
		HeaderLen  [4]bit.Bit
		Reserved   [3]bit.Bit `bit:"skip"`
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
		fmt.Printf("error:%s", err)
	}

	fmt.Printf("src=%d dst=%d\n", s.SrcPort, s.DstPort)
	fmt.Printf("SeqNo=%d AckNo=%d\n", s.SeqNo, s.AckNo)
	fmt.Printf("HeaderLen(raw)=%v\n", s.HeaderLen)
	fmt.Printf("Ack=%t\n", s.ACK)
}
```

## Struct Tag

The package supports struct tags.

|Tag|Description|
|---|-----------|
|`` `bit:"skip"` ``|Ignore the field. Offset is updated by the size of the field. It is useful for reserved field.|
|`` `bit:"-"` `` |Ignore the field. Offset is not updated.|
|`` `bit:"BE"` ``|Decode the field as big endian. It is useful for mixed endian data.|
|`` `bit:"LE"` ``|Decode the field as little endian. It is useful for mixed endian data.|

## Tool
* [readbit](v2/cmd/readbit/README.md)

## Document

https://pkg.go.dev/github.com/nokute78/go-bit/v2

## Old Document

https://godoc.org/github.com/nokute78/go-bit/pkg/bit

## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)