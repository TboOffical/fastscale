package utils

import "encoding/hex"

func HexStringToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}

func If(condition bool, true interface{}, false interface{}) interface{} {
	if condition {
		return true
	} else {
		return false
	}
}

func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
