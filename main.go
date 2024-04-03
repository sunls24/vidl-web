package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"vidl-web/common/logger"
	"vidl-web/router"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	r := gin.New()
	r.Use(logger.Middleware(), gin.Recovery())
	router.SetRouter(r, webFS)
	if err := r.Run(); err != nil {
		logger.Fatal().Err(err).Msg("server run failed")
	}
}
