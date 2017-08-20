package packages

import (
	"fmt"
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/codegangsta/cli"
)

func CmdListPackages(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	//document.GetFile(tag)

	for _, p := range DistinctPackageList(tag, c.GlobalBool("force")) {
		fmt.Println(p)
	}
}
