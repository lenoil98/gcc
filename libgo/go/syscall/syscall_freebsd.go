// Copyright 2009,2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd
// +build freebsd

package syscall

import (
	"sync"
	"unsafe"
)

const (
	_SYS_FSTAT_FREEBSD12         = 551 // { int fstat(int fd, _Out_ struct stat *sb); }
	_SYS_FSTATAT_FREEBSD12       = 552 // { int fstatat(int fd, _In_z_ char *path, _Out_ struct stat *buf, int flag); }
	_SYS_GETDIRENTRIES_FREEBSD12 = 554 // { ssize_t getdirentries(int fd, _Out_writes_bytes_(count) char *buf, size_t count, _Out_ off_t *basep); }
	_SYS_STATFS_FREEBSD12        = 555 // { int statfs(_In_z_ char *path, _Out_ struct statfs *buf); }
	_SYS_FSTATFS_FREEBSD12       = 556 // { int fstatfs(int fd, _Out_ struct statfs *buf); }
	_SYS_GETFSSTAT_FREEBSD12     = 557 // { int getfsstat(_Out_writes_bytes_opt_(bufsize) struct statfs *buf, long bufsize, int mode); }
	_SYS_MKNODAT_FREEBSD12       = 559 // { int mknodat(int fd, _In_z_ char *path, mode_t mode, dev_t dev); }
)

// See https://www.freebsd.org/doc/en_US.ISO8859-1/books/porters-handbook/versions.html.
var (
	osreldateOnce sync.Once
	osreldate     uint32
)

// INO64_FIRST from /usr/src/lib/libc/sys/compat-ino64.h
const _ino64First = 1200031

func supportsABI(ver uint32) bool {
	osreldateOnce.Do(func() { osreldate, _ = SysctlUint32("kern.osreldate") })
	return osreldate >= ver
}

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Fileno), unsafe.Sizeof(Dirent{}.Fileno))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Reclen), unsafe.Sizeof(Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Namlen), unsafe.Sizeof(Dirent{}.Namlen))
}

func Getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r0, _, e1 := Syscall(SYS_GETDIRENTRIES, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// nametomib converts a sysctl string name to its MIB representation.
func nametomib(name string) (mib []_C_int, err error) {
	// Query sysctl.name2oid {CTL_SYSCTL, CTL_SYSCTL_NAME2OID} = {0, 3} to convert string to MIB OID slice
	oid := []_C_int{_CTL_SYSCTL, _CTL_SYSCTL_NAME2OID}
	
	// Convert name string to bytes
	buf := []byte(name)
	
	// Determine output MIB buffer size
	var lenNeeded uintptr = 32 * unsafe.Sizeof(_C_int(0))
	mibBuf := make([]_C_int, 32)
	
	err = sysctl(oid, (*byte)(unsafe.Pointer(&mibBuf[0])), &lenNeeded, (*byte)(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if err != nil {
		return nil, err
	}
	
	count := lenNeeded / unsafe.Sizeof(_C_int(0))
	return mibBuf[:count], nil
}
