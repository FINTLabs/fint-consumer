package utils

import (
	"fmt"
	"os"
)

func GetTempDirectory() string {
	dir := ".temp"
	err := os.MkdirAll(dir, 0777)

	if err != nil {
		fmt.Println("Unable to create .temp")
		os.Exit(2)
	}

	return dir
}
