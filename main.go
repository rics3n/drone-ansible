package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "ansible plugin"
	app.Usage = "ansible plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "inventory-path",
			Usage:  "ansible inventory path",
			Value:  defaultInventoryPath,
			EnvVar: "PLUGIN_INVENTORY_PATH",
		},
		cli.StringSliceFlag{
			Name:   "inventories",
			Usage:  "ansible inventory files",
			Value:  &cli.StringSlice{"staging"},
			EnvVar: "PLUGIN_INVENTORY,PLUGIN_INVENTORIES",
		},
		cli.StringFlag{
			Name:   "playbook",
			Usage:  "ansible playbook to execute",
			Value:  defaultPlaybook,
			EnvVar: "PLUGIN_PLAYBOOK",
		},
		cli.StringFlag{
			Name:   "ssh-key",
			Usage:  "ssh key to access remote hosts",
			EnvVar: "SSH_KEY,PLUGIN_SSH_KEY",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "project base path",
			EnvVar: "DRONE_WORKSPACE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Build: Build{
			Path: c.String("path"),
		},
		Config: Config{
			InventoryPath: c.String("inventory-path"),
			Inventories:   c.StringSlice("inventories"),
			Playbook:      c.String("playbook"),
			SSHKey:        c.String("ssh-key"),
		},
	}

	return plugin.Exec()
}
