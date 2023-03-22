package main

import (
	"net/http"

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
	logrus.Debugf("srv: whitelisting %v", list)
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(list...),
	}
	if cfg.CA.Email != "" {
		m.Email = cfg.CA.Email
	}
	if cfg.CA.CacheDir != "" {
		m.Cache = autocert.DirCache(cfg.CA.CacheDir)
	}
	// redirect http traffic to https
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("got request on http: %s %s %s", r.Host, r.Method, r.RequestURI)
			w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
			http.Redirect(w, r, "https://"+r.Host, http.StatusMovedPermanently)
		})
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			panic(err)
		}
	}()
	return cli.NewExitError(http.Serve(m.Listener(), rph), 2)
}
