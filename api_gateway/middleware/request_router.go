package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReverseProxyMiddleware(target string, stripPattern string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		c.Request.URL.Path = strings.Replace(c.Request.URL.Path, stripPattern, "", 1)
		c.Request.Host = targetURL.Host
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func ProxySwaggerDoc(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := http.Get(target)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Cannot fetch swagger doc"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}
