package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mattfarina/helm-repo-audit/audit"
	"github.com/mitchellh/go-homedir"

	// using urfave instead of github.com/spf13/cobra because this one has
	// releases.
	"gopkg.in/urfave/cli.v1"
)

// The configuration loaded. It would be super awesome if the context had some
// form of dependency injection containers... but alas it does not
var cfg audit.Configs

func main() {

	// The default location for the cache
	var cacheLoc string
	if ho, err := homedir.Dir(); err == nil {
		cacheLoc = filepath.Join(ho, ".config", "helm-repo-audit")
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			// uh oh
			panic(fmt.Sprintf("Unable to handle cache locations"))
		}
		cacheLoc = filepath.Join(cwd, "helm-repo-audit")
	}

	app := cli.NewApp()
	app.Name = "Helm Repo Audit"
	app.Usage = "Audit a Helm repository"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Location of the config file (required)",
			EnvVar: "HELM_REPO_AUDIT_CONFIG",
		},
		cli.StringFlag{
			Name:   "store",
			Usage:  "Location of persistent store",
			Value:  cacheLoc,
			EnvVar: "HELM_REPO_AUDIT_STORE",
		},

		// TODO(mattfarina): Add verbose option
	}

	// Load configuration prior to doing anything else
	app.Before = func(c *cli.Context) error {

		// Load config
		pth := c.GlobalString("config")
		if pth == "" {
			return errors.New("configuration not passed in")
		}

		var err error
		cfg, err = audit.LoadConfigFromFile(pth)
		if err != nil {
			return err
		}
		return nil
	}

	app.Commands = []cli.Command{
		// TODO(mattfarina): Add command to take in input and generate config
		// TODO(mattfarina): Add command to validate config files

		{
			Name:    "audit",
			Aliases: []string{"a"},
			Usage:   "Audit one or more Helm repositories",
			Action: func(c *cli.Context) error {
				return audit.Audit(cfg, c.GlobalString("store"))
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
