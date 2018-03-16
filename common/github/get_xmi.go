package github

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/FINTprosjektet/fint-consumer/common/utils"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/net/context"
	"golang.org/x/text/encoding/charmap"
)

func GetXMIFile(owner string, repo string, tag string, filename string, force bool) string {
	outFileName := getFilePath(tag)

	if force {
		downloadFile(owner, repo, tag, filename, outFileName)
		cleanFile(outFileName)
	} else if !utils.FileExists(outFileName) {
		downloadFile(owner, repo, tag, filename, outFileName)
		cleanFile(outFileName)
	}

	return outFileName
}

func downloadFile(owner string, repo string, tag string, filename string, outFileName string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{
		Ref: tag,
	}
	out, err := client.Repositories.DownloadContents(ctx, owner, repo, filename, opt)
	if err != nil {
		fmt.Printf("Unable to download XMI file from GitHub: %s", err)
	}
	outFile, err := os.Create(outFileName)
	defer outFile.Close()
	_, err = io.Copy(outFile, out)
	if err != nil {
		fmt.Printf("Unable to write XMI file: %s", err)
	}
}

func getFilePath(tag string) string {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to get homedir.")
		os.Exit(2)
	}
	dir := fmt.Sprintf("%s/.fint-consumer/.cache", homeDir)
	err = os.MkdirAll(dir, 0777)

	if err != nil {
		fmt.Println("Unable to create .fint-consumer")
		os.Exit(2)
	}

	outFileName := fmt.Sprintf("%s/%s.xml", dir, tag)

	return outFileName
}

func cleanFile(fileName string) {

	toUtf8(fileName)

	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string

	newLines = append(newLines, "<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	keep := false
	for i, line := range lines {
		if strings.Contains(line, "<xmi:Extension extender=\"Enterprise Architect\" extenderID=\"6.5\">") {
			keep = true
		}
		if strings.Contains(line, "</xmi:XMI>") {
			keep = false
		}
		if keep {
			l := strings.Replace(lines[i], "uml:", "", -1)
			l = strings.Replace(l, "xmi:", "", -1)
			newLines = append(newLines, l)
		}
	}

	output := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0777)
	if err != nil {
		log.Fatalln(err)
	}

}

func toUtf8(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening %s (%s)", fileName, err)
		os.Exit(2)
	}
	defer f.Close()

	r := charmap.Windows1252.NewDecoder().Reader(f)

	content, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, content, 0777)

	if err != nil {
		fmt.Println("\nio.Copy failed:", err)
	}

}
