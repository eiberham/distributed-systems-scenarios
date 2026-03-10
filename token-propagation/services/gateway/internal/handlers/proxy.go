package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v5"
)

func Proxy(target string) echo.HandlerFunc {
	targetURL, _ := url.Parse(target)

	return func(c *echo.Context) error {
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Keep original path, only override if target has a specific path
			if targetURL.Path != "" && targetURL.Path != "/" {
				req.URL.Path = targetURL.Path
			}
			req.Host = targetURL.Host
		}

		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
