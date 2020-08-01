// +build darwin linux
package mmap

import (
	"syscall"
	"unsafe"
)

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
