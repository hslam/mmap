// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin linux

package mmap

import (
	"syscall"
	"unsafe"
)

const (
	PROT_READ  = syscall.PROT_READ
	PROT_WRITE = syscall.PROT_WRITE
	PROT_EXEC  = syscall.PROT_EXEC

	MAP_SHARED  = syscall.MAP_SHARED
	MAP_PRIVATE = syscall.MAP_PRIVATE
	MAP_COPY    = MAP_PRIVATE
)

func protFlags(p PROT) (prot int, flags int) {
	prot = PROT_READ
	flags = MAP_SHARED
	if p&WRITE != 0 {
		prot |= PROT_WRITE
	}
	if p&COPY != 0 {
		flags = MAP_COPY
	}
	if p&EXEC != 0 {
		prot |= PROT_EXEC
	}
	return
}

func mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return syscall.Mmap(fd, offset, length, prot, flags)
}

func msync(b []byte) (err error) {
	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), syscall.MS_SYNC)
	if errno != 0 {
		err = syscall.Errno(errno)
	}
	return
}

func munmap(b []byte) (err error) {
	return syscall.Munmap(b)
}
