package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/elauqsap/echo-postgres-json-api/api"
	"github.com/elauqsap/echo-postgres-json-api/database"
	"github.com/elauqsap/echo-postgres-json-api/reverse"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"

	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

// config contains the global settings for the server and database
type config struct {
	Server   api.Config      `json:"server"`
	Database database.Config `json:"database"`
}

var conf config

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
			Name:   "run",
			Usage:  "start the server",
			Before: load,
			Action: func(c *cli.Context) error {
				store, _ := conf.Database.New()
				dfile, err := os.OpenFile("./db.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
				if err != nil {
					return err
				}
				sfile, err := os.OpenFile("./app.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
				if err != nil {
					return err
				}
				// configure and run the server here
				server, err := conf.Server.New(sfile, dfile, store)
				if err != nil {
					return err
				}
				server.Run(standard.WithConfig(engine.Config{
					Address: conf.Server.Bind,
				}))
				return nil
			},
		},
		{
			Name:  "keygen",
			Usage: "generate encrypted keys",
			Action: func(c *cli.Context) error {
				fmt.Print("Enter Password: ")
				pass, _ := terminal.ReadPassword(0)
				key := make([]byte, 32)
				rand.Read(key)
				if enc, err := reverse.Encrypt(key, pass); err == nil {
					fmt.Println("[+] Key: " + enc.Key)
					fmt.Println("[+] Cipher: " + enc.Cipher)
				} else {
					return errors.New("[-] " + err.Error())
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}

// loads the configurations into a structure
// to be used as a global operator at runtime
func load(c *cli.Context) error {
	if c.GlobalIsSet("config") {
		data, _ := ioutil.ReadFile(c.GlobalString("config"))
		if err := json.Unmarshal(data, &conf); err != nil {
			return errors.New("--config is required, check JSON format and file path")
		}
	}
	return nil
}
