package node

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/daweth/gevm/gevmtypes"
	"github.com/ethereum/go-ethereum/rlp"
)

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

// RawTxToTxObject converts a []bytes transaction to a transaction struct
func RawTxToTxObject(rawTxHex string) gevmtypes.Transaction {
	// example
	// rawTxHex := "0xf86b808502540be4008252089411111111111111111111111111111111111111111880de0b6b3a7640000801ba0b3ef7e65d0516b4a5b1db67e49b1d14d9c3e44e6fe288997560e6099cfec5c17a07e094d2502257f2d3a8d4eeb2e2d5a7170b2f1c5b6d8e6f6a5a8d87c465dd771b3a"

	// Strip the 0x prefix if present
	if rawTxHex[:2] == "0x" {
		rawTxHex = rawTxHex[2:]
	}

	// Decode the hex string to a byte slice
	rlpBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		log.Fatalf("Failed to decode hex string: %v", err)
	}

	// RLP decode the byte slice back into a Transaction
	var tx gevmtypes.Transaction
	err = rlp.DecodeBytes(rlpBytes, &tx)
	if err != nil {
		log.Fatalf("Failed to RLP decode transaction: %v", err)
	}

	fmt.Printf("Decoded TX: %+v\n", tx)
	return tx
}

// check that an interface is a string
func interfaceIsString(i interface{}) bool {
	_, ok := i.(string)
	return ok
}
