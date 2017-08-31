package generate

import (
	"github.com/codegangsta/cli"
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/github"
)

func CmdGenerate(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	Generate(tag, force)

}
