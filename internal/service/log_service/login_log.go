package log_service

import (
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/enum"
	"github.com/gin-gonic/gin"
)

// NewLoginSuccess 登录成功
func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	//TODO:通过jwt获取username
	//token := c.GetHeader("token")
	UserID := uint(1)
	username := " "
	global.DB.Create(&model.LogModel{
		Type:        enum.LogLoginType,
		Title:       "用户登录",
		Content:     "",
		UserID:      UserID,
		IP:          ip,
		LoginStatus: true,
		Username:    username,
		Password:    "-",
		LoginType:   loginType,
	})
}

// NewLoginFail 登录失败
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, password string) {
	ip := c.ClientIP()
	//登录失败无用户id
	global.DB.Create(&model.LogModel{
		Type:        enum.LogLoginType,
		Title:       "用户登录失败",
		Content:     msg,
		IP:          ip,
		LoginStatus: false,
		Username:    username,
		Password:    password,
		LoginType:   loginType,
	})
}
