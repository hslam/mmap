// +build !darwin,!linux,!windows

package mmap

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

func protFlags(p PROT) (prot int, flags int) {
	return 0, 0
}

type mmapper struct {
	sync.Mutex
	active map[*byte]*f
}

type f struct {
	fd  int
	buf []byte
}

func (m *mmapper) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	if length <= 0 {
		return nil, syscall.EINVAL
	}
	buf := make([]byte, length)
	cursor, _ := syscall.Seek(fd, 0, os.SEEK_CUR)
	syscall.Seek(fd, 0, os.SEEK_SET)
	n, err := syscall.Read(fd, buf)
	syscall.Seek(fd, cursor, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	if n < length {
		return nil, errors.New("length > file size")
	}
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&buf[0])), length, length}
	b := *(*[]byte)(unsafe.Pointer(&sl))
	p := &b[cap(b)-1]
	m.Lock()
	defer m.Unlock()
	m.active[p] = &f{fd, b}
	return b, nil
}

func (m *mmapper) Msync(b []byte) (err error) {
	if len(b) == 0 || len(b) != cap(b) {
		return syscall.EINVAL
	}
	p := &b[cap(b)-1]
	m.Lock()
	defer m.Unlock()
	f := m.active[p]
	if f.buf == nil || &f.buf[0] != &b[0] {
		return syscall.EINVAL
	}
	cursor, _ := syscall.Seek(f.fd, 0, os.SEEK_CUR)
	syscall.Seek(f.fd, 0, os.SEEK_SET)
	_, err = syscall.Write(f.fd, b)
	syscall.Seek(f.fd, cursor, os.SEEK_SET)
	return err
}

func (m *mmapper) Munmap(data []byte) (err error) {
	if len(data) == 0 || len(data) != cap(data) {
		return syscall.EINVAL
	}
	p := &data[cap(data)-1]
	m.Lock()
	defer m.Unlock()
	f := m.active[p]
	if f.buf == nil || &f.buf[0] != &data[0] {
		return syscall.EINVAL
	}
	cursor, _ := syscall.Seek(f.fd, 0, os.SEEK_CUR)
	syscall.Seek(f.fd, 0, os.SEEK_SET)
	_, err = syscall.Write(f.fd, data)
	syscall.Seek(f.fd, cursor, os.SEEK_SET)
	delete(m.active, p)
	f = nil
	return err
}

var mapper = &mmapper{
	active: make(map[*byte]*f),
}

func mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func msync(b []byte) (err error) {
	return mapper.Msync(b)
}

func munmap(b []byte) (err error) {
	return mapper.Munmap(b)
}
