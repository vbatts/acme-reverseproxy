package main

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vbatts/acme-reverseproxy/proxymap"
	"golang.org/x/crypto/acme/autocert"
)

func srvCommand(c *cli.Context) error {
	list := []string{}
	for key := range cfg.Mapping {
		if key != "" {
			list = append(list, key)
		}
	}
	rpm, err := proxymap.ToReverseProxyMap(cfg.Mapping)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	rph := proxymap.NewReverseProxiesHandler(rpm)
	logrus.Debugf("srv: whitelisting %q", strings.Join(list, ","))
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(strings.Join(list, ",")),
	}
	if cfg.CA.Email != "" {
		m.Email = cfg.CA.Email
	}
	if cfg.CA.CacheDir != "" {
		m.Cache = autocert.DirCache(cfg.CA.CacheDir)
	}
	return cli.NewExitError(http.Serve(autocert.NewListener(list...), rph), 2)
}
