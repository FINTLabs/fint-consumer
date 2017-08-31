package setup

import (
	"log"
	"strings"
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
	"strconv"
)

func reportNeedOfChanges(name string) {

	fmt.Println("\n\nStart searching for change needs!")

	rootPath := fmt.Sprintf("%s/", getConsumerName(name))

	filepath.Walk(rootPath, needChange)
	fmt.Println("\n\nFinished searching for change needs!")
}

func needChange(path string, info os.FileInfo, err error) error {

	if info.IsDir() {
		return nil
	}
	report := ""
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "***fixme***") {
			report += fmt.Sprintf("%s. %s\n", strconv.Itoa(i+1), line)
		}
	}

	if len(report) > 0 {
		fmt.Printf("\n\nSearching file %s\n", path)
		fmt.Printf(report)
	}

	return nil
}
