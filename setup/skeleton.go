package setup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/FINTLabs/fint-consumer/common/config"
	"github.com/FINTLabs/fint-consumer/common/github"
	"github.com/FINTLabs/fint-consumer/common/utils"
)

func setupSkeleton(name string, reference string) {

	consumerName := getConsumerName(name)
	shortConsumerName := strings.Replace(consumerName, "fint-", "", -1)

	fmt.Printf("Creating %s ...\n", consumerName)
	fmt.Println("  > Cleanup old project files ...")

	os.RemoveAll(utils.GetWorkingDir(consumerName))

	fmt.Printf("  > Cloning repository %s, %s ...\n", config.CONSUMER_SKELETON_URL, reference)

	err := github.Clone(consumerName, config.CONSUMER_SKELETON_URL, reference)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("--> Cleanup project ...")
	os.RemoveAll(utils.GetDotGitDir(consumerName))

	fmt.Println("--> Renaming project ...")
	err = filepath.Walk(utils.GetWorkingDir(consumerName), func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !!fi.IsDir() {
			return nil
		}

		read, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), "consumer-skeleton", shortConsumerName, -1)

		return ioutil.WriteFile(path, []byte(newContents), 0)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("--> Finished!")

}
func getConsumerName(name string) string {
	return fmt.Sprintf("fint-consumer-%s", name)
}
