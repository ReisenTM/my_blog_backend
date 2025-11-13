package routers

import (
	"blog/internal/api"
	"blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouters(gr *gin.RouterGroup) {
	app := api.App.UserApi
	gr.POST("auth/email-code", app.SendEmail)
	gr.POST("auth/register", middleware.EmailVerifyMiddleware, app.EmailRegister)
}
