# mmap
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
	length := len(str)
	file.Truncate(int64(length))
	file.Sync()
	prot, flags := mmap.ProtFlags(mmap.READ | mmap.WRITE)
	b, err := mmap.Mmap(int(file.Fd()), 0, length, prot, flags)
	if err != nil {
		panic(err)
	}
	copy(b, []byte(str))
	mmap.Msync(b)
	mmap.Munmap(b)
	file.Sync()
	buf := make([]byte, length)
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


### Authors
mmap was written by Meng Huang.


