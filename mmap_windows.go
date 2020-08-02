// +build windows

package mmap

import (
	"sync"
	"syscall"
	"unsafe"
)

const (
	PAGE_READONLY          = syscall.PAGE_READONLY
	PAGE_READWRITE         = syscall.PAGE_READWRITE
	PAGE_WRITECOPY         = syscall.PAGE_WRITECOPY
	PAGE_EXECUTE_READ      = syscall.PAGE_EXECUTE_READ
	PAGE_EXECUTE_READWRITE = syscall.PAGE_EXECUTE_READWRITE
	PAGE_EXECUTE_WRITECOPY = syscall.PAGE_EXECUTE_WRITECOPY

	FILE_MAP_COPY    = syscall.FILE_MAP_COPY
	FILE_MAP_WRITE   = syscall.FILE_MAP_WRITE
	FILE_MAP_READ    = syscall.FILE_MAP_READ
	FILE_MAP_EXECUTE = syscall.FILE_MAP_EXECUTE
)

type mmapper struct {
	sync.Mutex
	active map[*byte][]byte
}

func (m *mmapper) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	if length <= 0 {
		return nil, syscall.EINVAL
	}
	handle, err := syscall.CreateFileMapping(syscall.Handle(fd), nil, uint32(prot), 0, uint32(length), nil)
	if err != nil {
		return nil, err
	}

	addr, err := syscall.MapViewOfFile(handle, uint32(flags), 0, 0, uintptr(length))
	if err != nil {
		return nil, err
	}
	err = syscall.CloseHandle(syscall.Handle(handle))
	if err != nil {
		return nil, err
	}
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{addr, length, length}
	b := *(*[]byte)(unsafe.Pointer(&sl))
	p := &b[cap(b)-1]
	m.Lock()
	defer m.Unlock()
	m.active[p] = b
	return b, nil
}

func (m *mmapper) Msync(b []byte) (err error) {
	slice := (*struct {
		addr uintptr
		len  int
		cap  int
	})(unsafe.Pointer(&b))
	return syscall.FlushViewOfFile(slice.addr, uintptr(slice.len))
}

func (m *mmapper) Munmap(data []byte) (err error) {
	if len(data) == 0 || len(data) != cap(data) {
		return syscall.EINVAL
	}
	p := &data[cap(data)-1]
	m.Lock()
	defer m.Unlock()
	b := m.active[p]
	if b == nil || &b[0] != &data[0] {
		return syscall.EINVAL
	}
	err = syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&b[0])))
	if err != nil {
		return err
	}
	delete(m.active, p)
	return nil
}

var mapper = &mmapper{
	active: make(map[*byte][]byte),
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