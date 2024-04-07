package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"vidlp/common/utils"
)

func Proxy(c *gin.Context) {
	target, exit := utils.MustQuery(c, "url")
	if exit {
		return
	}
	req, _ := http.NewRequest(c.Request.Method, target, c.Request.Body)
	for k, v := range c.Request.Header {
		if strings.ToLower(k) == "referer" {
			continue
		}
		req.Header[k] = v
	}

	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	_, _ = io.Copy(c.Writer, resp.Body)
}
