package user_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/utils/jwts"
	"blog/internal/utils/pwd"
	"github.com/gin-gonic/gin"
)

type EmailLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (UserApi) EmailLogin(c *gin.Context) {
	var req EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.FailWithMsg("请求体结构错误", c)
		return
	}
	var user model.UserModel
	err := global.DB.Take(&user, "email = ?", req.Email).Error
	if err != nil {
		resp.FailWithMsg("邮箱未注册", c)
		return
	}
	//校验密码
	if ok := pwd.CompareHashAndPassword(user.Password, req.Password); ok != true {
		resp.FailWithMsg("密码错误", c)
		return
	}
	//颁发token
	claims := jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	token, err := jwts.GetToken(claims)
	if err != nil {
		resp.FailWithMsg("生成token失败", c)
		return
	}
	resp.OkWithData(token, c)
}
