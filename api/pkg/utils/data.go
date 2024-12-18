package utils

import (
	"encoding/binary"
	"math/big"
)

func BytesToUint64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func Uint64ToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BigIntToBytes(i *big.Int) []byte {
	return i.Bytes()
}

func BytesEqual(a, b []byte) bool {
	maxLength := len(a)
	if len(b) > maxLength {
		maxLength = len(b)
	}

	for i := 0; i < maxLength; i++ {
		byteA := byte(0)
		byteB := byte(0)

		if i < len(a) {
			byteA = a[i]
		}
		if i < len(b) {
			byteB = b[i]
		}

		if byteA != byteB {
			return false
		}
	}
	return true
}
