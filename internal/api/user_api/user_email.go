package user_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/enum"
	"blog/internal/service/email_service"
	"blog/internal/utils/email_store"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

const (
	EmailRegType = iota + 1
	EmailResetType
)

type SendEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Type  uint8  `json:"type" binding:"oneof=1 2"` //1注册 2 重置密码
}

type SendEmailResponse struct {
	CodeID string `json:"codeID"` //验证码id
}

// SendEmail 发送操作提示邮件
func (UserApi) SendEmail(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.FailWithMsg("请求体格式错误", c)
		return
	}

	//生成验证码和id.存到验证码存储器中
	code := base64Captcha.RandText(4, "1234567890")
	id := base64Captcha.RandomId()
	var err error
	switch req.Type {
	case EmailRegType:
		//先查邮箱是否存在
		var um model.UserModel
		err = global.DB.Take(&um, "email = ?", req.Email).Error
		if err == nil && um.Email != "" {
			resp.FailWithMsg("邮箱已存在，请登录", c)
			return
		}
		err = email_service.SendRegCode(req.Email, code)
	case EmailResetType:
		var user model.UserModel
		err = global.DB.Take(&user, "email = ?", req.Email).Error
		if err != nil {
			resp.FailWithMsg("该邮箱不存在", c)
			return
		}
		// 还必须得是邮箱注册的或者已经绑定邮箱的
		if !(user.RegSource == enum.RegEmailSourceType || user.Email != "") {
			resp.FailWithMsg("仅支持邮箱注册或绑定邮箱的用户重置密码", c)
			return
		}
		err = email_service.SendResetCode(req.Email, code)
	}
	if err != nil {
		logrus.Errorf("邮件发送失败 %s", err)
		resp.FailWithMsg("发送邮件失败,%s", c)
		return
	}
	//err = global.CaptchaStore.Set(id, code)
	email_store.Set(id, req.Email, code)

	resp.OKWithMsg("验证码已发送", c)
}
