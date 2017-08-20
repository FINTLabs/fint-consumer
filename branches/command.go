package branches

import (
	"fmt"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/codegangsta/cli"
)

func CmdListBranches(c *cli.Context) {
	for _, b := range github.GetBranchList() {
		fmt.Println(b)
	}
}
