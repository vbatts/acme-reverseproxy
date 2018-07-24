package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vbatts/acme-reverseproxy/config"
)

var cfg config.Config

func main() {
	app := cli.NewApp()
	app.Name = "acme-reverseproxy"
	app.Usage = "A TLS-serving reverse-proxy, with the certificates generated from LetsEncrypt"
	app.Authors = []cli.Author{
		{Name: "Vincent Batts", Email: "vbatts@hashbangbash.com"},
	}
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:        "gen",
			Description: "generators of sorts",
			Subcommands: []cli.Command{
				{
					Name:        "config",
					Description: "generate a sample mapping configuration",
					Action:      genConfigAction,
				},
			},
		},
		{
			Name:        "srv",
			Description: "Start the reverseproxy server",
			Action:      srvCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Value: filepath.Join(os.Getenv("HOME"), ".acme-reverseproxy.toml"),
					Usage: "Configuration of mapping of hostname -> listener",
				},
			},
			Before: beforeAction,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func beforeAction(c *cli.Context) error {
	// find and read in the toml config file for hostname -> listener mappingj
	buf, err := ioutil.ReadFile(c.String("config"))
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Errorf("No config file found at %q. Try 'gen config'", c.String("config"))
		}
		return err
	}
	tmpConfig := config.Config{}
	if err := toml.Unmarshal(buf, &tmpConfig); err != nil {
		return err
	}
	cfg = tmpConfig
	return nil
}
