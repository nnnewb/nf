//go:build linux

package tuntap

import (
	"os"
	"syscall"
	"unsafe"

	/*
		#cgo CFLAGS: -Wall -Wpedantic
		#include <stdlib.h>

		int AllocateTun(char* dev);
		int GoErrno();
	*/
	"C"
)

func CAllocateTun(dev string) (*os.File, error) {
	cstr := C.CString(dev)
	defer C.free(unsafe.Pointer(cstr))

	fd := C.AllocateTun(cstr)
	if fd < 0 {
		return nil, syscall.Errno(C.GoErrno())
	}

	return os.NewFile(uintptr(fd), "/dev/net/tun"), nil
}
