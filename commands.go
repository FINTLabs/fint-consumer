package main

import (
	"fmt"
	"os"

	"github.com/FINTLabs/fint-consumer/branches"
	"github.com/FINTLabs/fint-consumer/generate"
	"github.com/FINTLabs/fint-consumer/packages"
	"github.com/FINTLabs/fint-consumer/setup"
	"github.com/FINTLabs/fint-consumer/tags"
	"github.com/codegangsta/cli"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "GITHUB_OWNER",
		Name:   "owner",
		Value:  "FINTLabs",
		Usage:  "Git repository containing model",
	},
	cli.StringFlag{
		EnvVar: "GITHUB_PROJECT",
		Name:   "repo",
		Value:  "fint-informasjonsmodell",
		Usage:  "Git repository containing model",
	},
	cli.StringFlag{
		EnvVar: "MODEL_FILENAME",
		Name:   "filename",
		Value:  "FINT-informasjonsmodell.xml",
		Usage:  "File name containing information model",
	},
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
		Flags:  []cli.Flag{},
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
	{
		Name:   "setup",
		Usage:  "setup a consumer project",
		Action: setup.CmdSetupConsumer,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "name of the consumer, e.g. personal",
			},
			cli.StringFlag{
				Name:  "component, c",
				Usage: "component prefix, e.g. administrasjon",
			},
			cli.StringFlag{
				Name:  "package, p",
				Usage: "the package you want to create the consumer for, e.g. kodeverk",
			},
			cli.BoolFlag{
				Name:  "includePerson",
				Usage: "Include person model",
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
