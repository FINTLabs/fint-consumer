package utils

import (
	"fmt"
	"os"
	"github.com/mitchellh/go-homedir"
)

func GetTempDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to get homedir.")
		os.Exit(2)
	}
	dir := fmt.Sprintf("%s/.fint-consumer/tmp", homeDir)
	err = os.MkdirAll(dir, 0777)

	if err != nil {
		fmt.Println("Unable to create .fint-consumer")
		os.Exit(2)
	}

	return dir
}
