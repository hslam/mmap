// +build windows

package mmap

import (
	"os"
	"testing"
)

func TestMmap(t *testing.T) {
	name := "mmap"
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(name)
	defer file.Close()
	size := 11
	file.Truncate(int64(size))
	file.Sync()
	b, err := Mmap(int(file.Fd()), 0, size, PAGE_READWRITE, FILE_MAP_WRITE)
	if err != nil {
		t.Error(err)
	}
	str := "Hello world"
	copy(b, []byte(str))
	if err := Msync(b); err != nil {
		t.Error(err)
	}
	buf := make([]byte, size)
	if _, err := file.Read(buf); err != nil {
		t.Error(err)
	}
	if str != string(buf) {
		t.Errorf("%s!=%s", str, string(buf))
	}
	if err := Munmap(b); err != nil {
		t.Error(err)
	}
	file.Sync()
}

func BenchmarkMmap(b *testing.B) {
	name := "mmap"
	file, err := os.Create(name)
	if err != nil {
		b.Error(err)
	}
	defer os.Remove(name)
	defer file.Close()
	size := 11
	file.Truncate(int64(size))
	file.Sync()
	d, err := Mmap(int(file.Fd()), 0, size, PAGE_READWRITE, FILE_MAP_WRITE)
	if err != nil {
		b.Error(err)
	}
	str := "Hello world"
	for i := 0; i < b.N; i++ {
		copy(d, []byte(str))
		Msync(d)
		file.Sync()
	}
	if err := Munmap(d); err != nil {
		b.Error(err)
	}
}
