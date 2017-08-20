package main

import (
	"fmt"
	"os"

	"github.com/FINTprosjektet/fint-consumer/branches"
	"github.com/FINTprosjektet/fint-consumer/classes"
	"github.com/FINTprosjektet/fint-consumer/generate"
	"github.com/FINTprosjektet/fint-consumer/namespaces"
	"github.com/FINTprosjektet/fint-consumer/packages"
	"github.com/FINTprosjektet/fint-consumer/tags"
	"github.com/codegangsta/cli"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "",
		Name:   "tag, t",
		Value:  "latest",
		Usage:  "the tag (version) of the model to generate",
	},
	cli.BoolFlag{
		EnvVar: "",
		Name:   "force, f",
		Usage:  "force downloading XMI for GitHub.",
	},
}

var Commands = []cli.Command{
	{
		Name:   "generate",
		Usage:  "generates consumer code",
		Action: generate.CmdGenerate,
		Flags: []cli.Flag{},
	},
	{
		Name:   "listPackages",
		Usage:  "list Java packages",
		Action: packages.CmdListPackages,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listTags",
		Usage:  "list tags",
		Action: tags.CmdListTags,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listBranches",
		Usage:  "list branches",
		Action: branches.CmdListBranches,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
