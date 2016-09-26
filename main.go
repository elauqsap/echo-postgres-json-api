package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/bgentry/speakeasy"
	"github.com/elauqsap/echo-postgres-json-api/api"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"

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
				if c.GlobalIsSet("config") {
					if Config = api.NewConfig(c.GlobalString("config")); Config == nil {
						return errors.New("--config is required, check JSON format and file path")
					}
				}
				// configure and run the server here
				server, err := Config.NewServer()
				if err != nil {
					return err
				}
				server.Run(standard.WithConfig(engine.Config{
					Address: Config.Server.Bind,
				}))
				return nil
			},
		},
		{
			Name:  "keygen",
			Usage: "generate encrypted keys",
			Action: func(c *cli.Context) error {
				plain, _ := speakeasy.Ask("[!] Encrypt: ")
				key := make([]byte, 32)
				rand.Read(key)
				if enc, err := encrypt(key, []byte(plain)); err == nil {
					fmt.Println("[+] Key: " + base64.StdEncoding.EncodeToString(key))
					fmt.Println("[+] Cipher: " + *enc)
				} else {
					return errors.New("[-] " + err.Error())
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
