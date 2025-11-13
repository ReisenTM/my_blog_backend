package user_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/utils/pwd"
	"github.com/gin-gonic/gin"
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
	newUser := &model.UserModel{
		Password: hashedPwd,
	}
	err = global.DB.Create(newUser).Error
	if err != nil {
		resp.FailWithMsg("用户创建失败", c)
		return
	}
	resp.OKWithMsg("注册成功", c)
}
