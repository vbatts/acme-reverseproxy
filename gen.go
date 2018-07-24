package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"github.com/vbatts/acme-reverseproxy/config"
)

func genConfigAction(c *cli.Context) error {
	tmpConfig := config.Config{
		CA: config.CA{
			Email:    "admin@example.com",
			CacheDir: "/tmp/acme-reverseproxy",
		},
		Mapping: map[string]string{
			"example.com": "http://localhost:5000",
		},
	}
	e := toml.NewEncoder(os.Stdout)
	if err := e.Encode(tmpConfig); err != nil {
		return err
	}
	return nil
}
