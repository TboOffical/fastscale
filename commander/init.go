package main

import (
	"errors"
	"os"
)

func checkDirectories() error {
	//Check if the data directory exists
	dir, err := os.ReadDir("./data")
	if err != nil {
		//If it doesn't exist, start the setup process
		return errors.New("setup required")
	}

	_ = dir
	return nil
}

func setupHandler() {

}
