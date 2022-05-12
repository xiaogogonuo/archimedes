// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package encryption

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/md4"
	"hash"
)

func ArchHashByte(src []byte, algorithm string) []byte {
	var _hash hash.Hash

	switch algorithm {
	case "md4":
		_hash = md4.New()
	case "sha1":
		_hash = sha1.New()
	case "sha256":
		_hash = sha256.New()
	case "sha512":
		_hash = sha512.New()
	default:
		_hash = md5.New()
	}

	_hash.Write(src)

	return _hash.Sum(nil)
}

func ArchHashString(src []byte, algorithm string) string {
	return fmt.Sprintf("%x", ArchHashByte(src, algorithm))
}
