package router

import (
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine, webFS embed.FS) {
	r.Use(static.Serve("/", static.EmbedFolder(webFS, "web/dist")))
	setApiRouter(r)
}

func setApiRouter(router *gin.Engine) {
	//r := router.Group("/api")
}
