package mmap

func mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return nil, nil
}

func msync(b []byte) (err error) {
	return
}

func munmap(b []byte) (err error) {
	return
}
