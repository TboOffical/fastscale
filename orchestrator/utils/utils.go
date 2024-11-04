package utils

import "encoding/hex"

func HexStringToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}
