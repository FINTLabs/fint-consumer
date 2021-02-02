package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetTempDirectory() string {
	dir, err := ioutil.TempDir("", "fint-consumer")

	if err != nil {
		fmt.Println("Unable to create tempdir fint-consumer")
		os.Exit(2)
	}

	return dir
}
