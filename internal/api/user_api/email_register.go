package user_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/enum"
	"blog/internal/utils/pwd"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type EmailRegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	EmailCode string `json:"emailCode" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (UserApi) EmailRegister(c *gin.Context) {
	var req EmailRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.FailWithError(err, c)
		return
	}
	hashedPwd, err := pwd.GenerateFromPassword(req.Password)
	if err != nil {
		resp.FailWithMsg(err.Error(), c)
		return
	}
	randName := "用户" + base64Captcha.RandText(4, "1234567890")
	newUser := &model.UserModel{
		Password:  hashedPwd,
		RegSource: EmailRegType,
		Role:      enum.RoleUserType,
		Email:     req.Email,
		Username:  randName,
	}
	err = global.DB.Create(newUser).Error
	if err != nil {
		resp.FailWithMsg("用户创建失败", c)
		return
	}
	resp.OKWithMsg("注册成功", c)
}
