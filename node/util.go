package node

import "encoding/binary"

// Uint64ToBytes converts a uint64 to a slice of 8 bytes
func Uint64ToBytes(num uint64) []byte {
	bytes := make([]byte, 8) // a uint64 always takes 8 bytes
	binary.BigEndian.PutUint64(bytes, num)
	return bytes
}
