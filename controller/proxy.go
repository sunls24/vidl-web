package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"vidl-web/common"
)

func Proxy(c *gin.Context) {
	targetUrl, exit := common.Query(c, "url")
	if exit {
		return
	}
	client := &http.Client{}
	req, _ := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	for k, v := range c.Request.Header {
		if strings.ToLower(k) == "referer" {
			continue
		}
		req.Header[k] = v
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorMsg(err.Error()))
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}

	_, _ = io.Copy(c.Writer, resp.Body)
}
