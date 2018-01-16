package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
	"github.com/FINTprosjektet/fint-consumer/common/config"
)

func setupSkeleton(name string) {

	consumerName := getConsumerName(name)
	shortConsumerName := strings.Replace(consumerName, "fint-", "", -1)

	fmt.Printf("Creating %s ...\n", consumerName)
	fmt.Println("  > Cleanup old project files ...")

	os.RemoveAll(utils.GetWorkingDir(consumerName))

	fmt.Printf("  > Cloning repository %s ...\n", config.CONSUMER_SKELETON_URL)

	err := github.Clone(consumerName, config.CONSUMER_SKELETON_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("--> Cleanup project ...")
	os.RemoveAll(utils.GetDotGitDir(consumerName))

	fmt.Println("--> Renaming project ...")
	err = filepath.Walk(utils.GetWorkingDir(consumerName), func (path string, fi os.FileInfo, err error) error {
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
