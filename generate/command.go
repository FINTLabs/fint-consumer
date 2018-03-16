package generate

import (
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/codegangsta/cli"
)

func CmdGenerate(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	Generate(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), force)

}
