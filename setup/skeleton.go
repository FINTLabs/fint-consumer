package setup

import (
	"fmt"
	"os"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
	"github.com/FINTprosjektet/fint-consumer/common/config"
)

func setupSkeleton(name string) {

	consumerName := getConsumerName(name)
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

	fmt.Println("--> Finished!")

}
func getConsumerName(name string) string {
	return fmt.Sprintf("fint-consumer-%s", name)
}
