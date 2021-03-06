package tags

import (
	"fmt"

	"github.com/FINTLabs/fint-consumer/common/github"
	"github.com/urfave/cli"
)

func CmdListTags(c *cli.Context) {
	for _, t := range github.GetTagList(c.GlobalString("owner"), c.GlobalString("repo")) {
		fmt.Println(t)
	}
}
