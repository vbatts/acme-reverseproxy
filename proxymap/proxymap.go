package proxymap

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// ToReverseProxyMap assembles a basic ReverseProxyMap from a map of domain name key, URL string value.
func ToReverseProxyMap(m map[string]string) (ReverseProxyMap, error) {
	rpm := ReverseProxyMap{}
	for k, v := range m {
		rpURL, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		rpm[k] = httputil.NewSingleHostReverseProxy(rpURL)
	}
	return rpm, nil
}

// ReverseProxyMap maps a domain name to the HTTP Handler for that domanin.
// Intended that the Handler is an httputil.ReverseProxy, it is flexible to other implementations of the HTTP Handler.
type ReverseProxyMap map[string]http.Handler

// NewReverseProxiesHandler produces a new HTTP handler for the provided set of reverse proxy mappings
func NewReverseProxiesHandler(rpm ReverseProxyMap) ReverseProxiesHandler {
	return &reverseProxiesHandler{
		Map:             rpm,
		NotFoundHandler: http.NotFoundHandler(),
	}
}

// ReverseProxiesHandler is an HTTP Handler for routing requests to the backing reverse proxies
type ReverseProxiesHandler interface {
	http.Handler
}

type reverseProxiesHandler struct {
	Map             ReverseProxyMap
	NotFoundHandler http.Handler
}

func (rph reverseProxiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		if addrErr, ok := err.(*net.AddrError); ok && !strings.Contains(addrErr.Err, "missing port") {
			logrus.Errorf("%T %#v", err, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		logrus.Debugf("using %q as hostname (no port)", r.Host)
		host = r.Host
	}
	logrus.Printf("request for %q %q %q", r.Host, r.URL, host)
	if v, ok := rph.Map[r.Host]; ok {
		v.ServeHTTP(w, r)
		return
	}

	if rph.NotFoundHandler == nil {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	rph.NotFoundHandler.ServeHTTP(w, r)
}
