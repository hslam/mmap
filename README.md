# mmap
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/mmap)](https://pkg.go.dev/github.com/hslam/mmap)
[![Build Status](https://travis-ci.org/hslam/mmap.svg?branch=master)](https://travis-ci.org/hslam/mmap)
[![codecov](https://codecov.io/gh/hslam/mmap/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/mmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/mmap?v=7e100)](https://goreportcard.com/report/github.com/hslam/mmap)
[![LICENSE](https://img.shields.io/github/license/hslam/mmap.svg?style=flat-square)](https://github.com/hslam/mmap/blob/master/LICENSE)

Package mmap provides a way to memory-map a file.

## Get started

### Install
```
go get github.com/hslam/mmap
```
### Import
```
import "github.com/hslam/mmap"
```
### Usage
#### Example
```
package main

import (
	"fmt"
	"github.com/hslam/mmap"
	"os"
)

func main() {
	name := "mmap"
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer os.Remove(name)
	defer file.Close()
	str := "Hello world"
	file.Truncate(int64(len(str)))
	b, err := mmap.Open(int(file.Fd()), 0, len(str), mmap.READ|mmap.WRITE)
	if err != nil {
		panic(err)
	}
	defer mmap.Munmap(b)
	copy(b, []byte(str))
	mmap.Msync(b)
	buf := make([]byte, len(str))
	if n, err := file.Read(buf); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(buf[:n]))
	}
}
```

### Output
```
Hello world
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
mmap was written by Meng Huang.


