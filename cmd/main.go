package main

import (
	"log"
	"os"

	go_signer "github.com/agglayer/go_signer"
	"github.com/agglayer/go_signer/cmd/version"
	cli "github.com/urfave/cli/v2"
)

const appName = "go_signer"

var ()

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = go_signer.Version

	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Application version and build",
			Action:  version.VersionCmd,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
