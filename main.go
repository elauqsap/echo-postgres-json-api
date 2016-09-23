package main

import (
	"errors"
	"os"

	"github.com/elauqsap/echo-postgres-json-api/api"

	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

// TODO: Implement Migrations
// TODO: Implement User Model
// TODO: Test Database
// TODO: Implement User API
// TODO: Test API
// IDEA: Blog About This =)

// Config contains the global settings for the server and database
var Config *api.Config

func main() {
	app := cli.NewApp()
	app.Name = "echo-postgres-json-api"
	app.Author = "Pasquale D'Agostino"
	app.Version = "1.0.0"
	app.Usage = "command line utilities for managing a server and backend"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "path to a valid config file (required)",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalIsSet("config") {
			if Config = api.NewConfig(c.GlobalString("config")); Config != nil {
				return nil
			}
		}
		return errors.New("--config is required, check JSON format and file path")
	}
	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() && c.NumFlags() < 1 {
			return cli.ShowAppHelp(c)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "start the server",
			Action: func(c *cli.Context) error {
				// configure and run the server here
				return nil
			},
		},
	}

	app.Run(os.Args)
}
