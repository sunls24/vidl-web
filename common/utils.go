package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Query(c *gin.Context, key string) (string, bool) {
	var value = c.Query(key)
	if value == "" {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("No %s parameters found", key))
		return value, true
	}
	return value, false
}

func GetEnv(key, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return value
}
