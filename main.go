package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "FINTProsjektet"
	app.Email = ""
	app.Usage = "Generates consumer code from EA XMI export. " +
		"This utility is mainly for internal FINT use, but if you " +
		"find it usefull, please use it!"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
