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
	offset := int64(os.Getpagesize() * 4)
	size := 11
	file.Truncate(int64(size) + offset)
	file.Sync()
	prot, flags := ProtFlags(READ | WRITE)
	m, err := Mmap(Fd(file), offset, size, prot, flags)
	if err != nil {
		t.Error(err)
	}
	str := "Hello world"
	copy(m, []byte(str))
	if err := Msync(m); err != nil {
		t.Error(err)
	}
	buf := make([]byte, size)
	file.Seek(offset, os.SEEK_SET)
	file.Read(buf)
	if str != string(buf) {
		t.Errorf("%s!=%s", str, string(buf))
	}
	if err := Munmap(m); err != nil {
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
	m, err := Open(Fd(file), 0, Fsize(file), READ|WRITE)
	if err != nil {
		b.Error(err)
	}
	str := "Hello world"
	for i := 0; i < b.N; i++ {
		copy(m, []byte(str))
		Msync(m)
		file.Sync()
	}
	if err := Munmap(m); err != nil {
		b.Error(err)
	}
}
