package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	host string
	port string
	path string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "api") || strings.Contains(r.URL.Path, "swagger") {
			next.ServeHTTP(w, r)
			return
		}
		target := &url.URL{
			Scheme: "http",
			Host:   rp.host + ":" + rp.port,
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Director = func(r *http.Request) {
			r.URL.Scheme = target.Scheme
			r.URL.Host = target.Host
			r.Host = target.Host
		}
		proxy.ServeHTTP(w, r)
	})
}
