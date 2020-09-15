// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package mmap provides a way to memory-map a file.
package mmap

import (
	"os"
)

// PROT is the prot.
type PROT int

const (
	// READ represents the read prot
	READ PROT = 1 << iota
	// WRITE represents the write prot
	WRITE
	// COPY represents the copy prot
	COPY
	// EXEC represents the exec prot
	EXEC
)

// Fd returns the integer file descriptor referencing the open file.
// The file descriptor is valid only until f.Close is called or f is garbage collected.
func Fd(f *os.File) int {
	return int(f.Fd())
}

// Fsize returns the file size.
func Fsize(f *os.File) int {
	cursor, _ := f.Seek(0, os.SEEK_CUR)
	ret, _ := f.Seek(0, os.SEEK_END)
	f.Seek(cursor, os.SEEK_SET)
	return int(ret)
}

// ProtFlags returns prot and flags by PROT p.
func ProtFlags(p PROT) (prot int, flags int) {
	return protFlags(p)
}

// Offset returns the valid offset.
func Offset(offset int64) int64 {
	pageSize := int64(os.Getpagesize())
	return offset / pageSize * pageSize
}

// Open opens a mmap
func Open(fd int, offset int64, length int, p PROT) (data []byte, err error) {
	prot, flags := protFlags(p)
	return mmap(fd, offset, length, prot, flags)
}

//Mmap calls the mmap system call.
func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mmap(fd, offset, length, prot, flags)
}

// Msync calls the msync system call.
func Msync(b []byte) (err error) {
	return msync(b)
}

// Munmap calls the munmap system call.
func Munmap(b []byte) (err error) {
	return munmap(b)
}
