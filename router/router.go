package router

import (
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"vidl-web/controller"
)

func SetRouter(r *gin.Engine, webFS embed.FS) {
	r.Use(static.Serve("/", static.EmbedFolder(webFS, "web/dist")))
	setApiRouter(r)
}

func setApiRouter(router *gin.Engine) {
	r := router.Group("/api")
	r.GET("/analyze", controller.Analyze)
	r.GET("/download", controller.Download)
	r.GET("/proxy", controller.Proxy)
}
