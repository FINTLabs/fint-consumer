package branches

import (
	"fmt"

	"github.com/FINTLabs/fint-consumer/common/github"
	"github.com/urfave/cli"
)

func CmdListBranches(c *cli.Context) {
	for _, b := range github.GetBranchList(c.GlobalString("owner"), c.GlobalString("repo")) {
		fmt.Println(b)
	}
}
