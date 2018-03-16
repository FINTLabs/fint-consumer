package tags

import (
	"fmt"

	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/codegangsta/cli"
)

func CmdListTags(c *cli.Context) {
	for _, t := range github.GetTagList(c.GlobalString("owner"), c.GlobalString("repo")) {
		fmt.Println(t)
	}
}
