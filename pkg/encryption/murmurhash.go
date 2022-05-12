// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package encryption

import "unsafe"

const m uint64 = 0xc6a4a7935bd1e995
const seed uint64 = 0xc70f6907

func MurmurHash64A1(src []byte) int64 {
	const r = 47
	var l = len(src)
	var h = seed ^ uint64(l)*m
	var data = src
	var l8 = l / 8
	var k uint64

	for i := 0; i < l8; i++ {
		i8 := i * 8
		k = *(*uint64)(unsafe.Pointer(&src[i8]))
		k *= m
		k ^= k >> r
		k *= m

		h ^= k
		h *= m
	}
	data = data[l8*8:]
	switch l & 7 {
	case 7:
		h ^= uint64(data[6]) << 48
		fallthrough
	case 6:
		h ^= uint64(data[5]) << 40
		fallthrough
	case 5:
		h ^= uint64(data[4]) << 32
		fallthrough
	case 4:
		h ^= uint64(data[3]) << 24
		fallthrough
	case 3:
		h ^= uint64(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint64(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint64(data[0])
		h *= m
	}
	h ^= h >> r
	h *= m
	h ^= h >> r
	return int64(h)
}

func MurmurHash64A2(src []byte) int64 {
	const r = 47
	var l = len(src)
	var h = seed ^ uint64(l)*m
	var data = src
	var l8 = l / 8
	var k uint64
	for i := 0; i < l8; i++ {
		i8 := i * 8
		k = uint64(data[i8+0]) + uint64(data[i8+1])<<8 +
			uint64(data[i8+2])<<16 + uint64(data[i8+3])<<24 +
			uint64(data[i8+4])<<32 + uint64(data[i8+5])<<40 +
			uint64(data[i8+6])<<48 + uint64(data[i8+7])<<56
		k *= m
		k ^= k >> r
		k *= m

		h ^= k
		h *= m
	}
	data = data[l8*8:]
	switch l & 7 {
	case 7:
		h ^= uint64(data[6]) << 48
		fallthrough
	case 6:
		h ^= uint64(data[5]) << 40
		fallthrough
	case 5:
		h ^= uint64(data[4]) << 32
		fallthrough
	case 4:
		h ^= uint64(data[3]) << 24
		fallthrough
	case 3:
		h ^= uint64(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint64(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint64(data[0])
		h *= m
	}
	h ^= h >> r
	h *= m
	h ^= h >> r
	return int64(h)
}
