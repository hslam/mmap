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
		fmt.Println(err)
	}
	defer os.Remove(name)
	defer file.Close()
	size := 11
	file.Truncate(int64(size))
	file.Sync()
	prot, flags := mmap.ProtFlags(mmap.READ | mmap.WRITE)
	b, err := mmap.Mmap(int(file.Fd()), 0, size, prot, flags)
	if err != nil {
		fmt.Println(err)
	}
	str := "Hello world"
	copy(b, []byte(str))
	if err := mmap.Msync(b); err != nil {
		fmt.Println(err)
	}
	file.Sync()
	buf := make([]byte, size)
	if n, err := file.Read(buf); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(buf[:n]))
	}
	if err := mmap.Munmap(b); err != nil {
		fmt.Println(err)
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


