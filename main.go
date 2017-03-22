package main

import (
	"net/url"
	"net/http"
	"net/http/httputil"
	"strings"
	"golang.org/x/crypto/acme/autocert"
	"crypto/tls"
	"fmt"
	"time"
	"log"
)

func main() {
	target := &url.URL{
		Scheme: "http",
		Host:   "www.amazingchatstories.com",
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}

		req.Header.Set("Cookie", "")
		req.Header.Set("User-Agent", fmt.Sprintf("%d", time.Now().Nanosecond()))

		req.Host = target.Host

	}

	server := &http.Server{
		Handler: proxy,
		Addr:    ":9090",
	}
	log.Println("open your browser on http://localhost:9090")
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
func redirectHttps(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func LETSENCRYPT() (*tls.Config) {
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("./certcache"),
	}
	return &tls.Config{
		GetCertificate: m.GetCertificate,
	}
}
