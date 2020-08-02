package mmap

import (
	"os"
)

func Fd(f *os.File) int {
	return int(f.Fd())
}

func Fsize(f *os.File) int {
	ret, _ := f.Seek(0, os.SEEK_END)
	return int(ret)
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
