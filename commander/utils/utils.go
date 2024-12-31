package utils

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/charmbracelet/log"
	"regexp"
)

func HexStringToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
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

func NullErr(value string) (string, error) {
	if len(value) == 0 {
		return "", errors.New("empty")
	}
	return value, nil
}

func IsInListStrings(object interface{}, list []string) bool {
	present := false

	for _, item := range list {
		if object == item {
			present = true
		}
	}

	return present
}

func IsInList(object interface{}, list []interface{}) bool {
	found := false
	for _, thing := range list {
		if thing == object {
			found = true
		}
	}
	return found
}

// TableColNameCheck makes sure that a given col name is okay to be put into a table True=Okay
func TableColNameCheck(name string) bool {
	//use regex
	matched, err := regexp.Match("/[^a-z_]/g", []byte(name))
	if err != nil {
		log.Error("An error occurred when checking if a col name was valid, regex: ", err)
		return false
	}
	return !matched
}
