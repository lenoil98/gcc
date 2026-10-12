// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gccgo does not implement the FreeBSD 11 dirent conversion.

//go:build !gccgo

package syscall

type Dirent_freebsd11 = dirent_freebsd11

var (
	Roundup              = roundup
	ConvertFromDirents11 = convertFromDirents11
)
