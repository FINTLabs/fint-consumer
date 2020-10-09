package generate

import (
	"github.com/FINTLabs/fint-consumer/common/config"
	"github.com/FINTLabs/fint-consumer/common/github"
	"github.com/urfave/cli"
)

func CmdGenerate(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	Generate(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), force, c.String("component"), c.String("package"), false)

}
