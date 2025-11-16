package routers

import (
	"blog/internal/api"
	"blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ArticleRouters(gr *gin.RouterGroup) {
	app := api.App.ArticleApi
	gr.GET("/posts", app.ArticleOverview)
	gr.GET("/post/:slug", app.ArticleDetail)
	gr.POST("/post", middleware.AuthMiddleware, app.PostArticle)
}
