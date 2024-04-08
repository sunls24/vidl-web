package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"vidlp/common/logger"
	"vidlp/common/utils"
	"vidlp/router"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	r := gin.New()
	r.Use(logger.Middleware(), gin.Recovery())
	router.SetRouter(r, webFS)

	var host = utils.GetEnv("HOST", "127.0.0.1")
	var port = utils.GetEnv("PORT", "3003")
	if err := r.Run(host + ":" + port); err != nil {
		logger.Fatal().Err(err).Msg("server run failed")
	}
}
