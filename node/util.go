package node

import "encoding/binary"

// Uint64ToBytes converts a uint64 to a slice of 8 bytes
func Uint64ToBytes(num uint64) []byte {
	bytes := make([]byte, 8) // a uint64 always takes 8 bytes
	binary.BigEndian.PutUint64(bytes, num)
	return bytes
}

// BytesToUint64 converts a slice of 8 bytes to a uint64
func BytesToUint64(bytes []byte) uint64 {
	if len(bytes) < 8 {
		panic("byte slice is too short to be converted to uint64")
	}
	return binary.BigEndian.Uint64(bytes)
}
