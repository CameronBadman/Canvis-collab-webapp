package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func CreateProxy(targetURL string) (gin.HandlerFunc, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path + req.URL.Path
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}

	return func(c *gin.Context) {
		log.Printf("Proxying request: %s %s", c.Request.Method, c.Request.URL.Path)
		proxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}
