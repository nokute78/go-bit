# go-bit

[![Build Status](https://travis-ci.org/nokute78/go-bit.svg?branch=master)](https://travis-ci.org/nokute78/go-bit)
[![Go Report Card](https://goreportcard.com/badge/github.com/nokute78/go-bit)](https://goreportcard.com/report/github.com/nokute78/go-bit)
[![GoDoc](https://godoc.org/github.com/nokute78/go-bit/pkg/bit?status.svg)](https://godoc.org/github.com/nokute78/go-bit/pkg/bit)

A library to read/write bits from a byte slice.

## Installation

```
$ go get github.com/nokute78/go-bit/pkg/bit
```

## Usage
```go
package main

import (
	"fmt"
	"github.com/nokute78/go-bit/pkg/bit"
)

func main() {
	b := []byte{0x78} /* 0111_1000 in bit */

	/* try to get 4bits(1111b) from 0111_1000 */
	off := bit.Offset{Byte: 0, Bit: 3}

	ret, err := bit.GetBits(b, off, 4)
	if err != nil {
		fmt.Printf("error:%s\n", err)
	}

	fmt.Printf("0x%x\n", ret) /* Print 0x0f = 1111b */
}
```
## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)