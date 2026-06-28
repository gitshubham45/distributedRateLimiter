package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ForwardRequest(c *gin.Context, target string) {
	targetUrl, err := url.Parse(target)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	c.Request.URL.Host = targetUrl.Host
	c.Request.URL.Scheme = targetUrl.Scheme
	c.Request.Host = targetUrl.Host

	proxy.ServeHTTP(c.Writer, c.Request)
}
