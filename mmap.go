// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package mmap provides a way to memory-map a file.
package mmap

import (
	"os"
)

type PROT int

const (
	READ PROT = 1 << iota
	WRITE
	COPY
	EXEC
)

func Fd(f *os.File) int {
	return int(f.Fd())
}

func Fsize(f *os.File) int {
	cursor, _ := f.Seek(0, os.SEEK_CUR)
	ret, _ := f.Seek(0, os.SEEK_END)
	f.Seek(cursor, os.SEEK_SET)
	return int(ret)
}

func ProtFlags(p PROT) (prot int, flags int) {
	return protFlags(p)
}

func Offset(offset int64) int64 {
	pageSize := int64(os.Getpagesize())
	return offset / pageSize * pageSize
}

func Open(fd int, offset int64, length int, p PROT) (data []byte, err error) {
	prot, flags := protFlags(p)
	return mmap(fd, offset, length, prot, flags)
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mmap(fd, offset, length, prot, flags)
}

func Msync(b []byte) (err error) {
	return msync(b)
}

func Munmap(b []byte) (err error) {
	return munmap(b)
}
