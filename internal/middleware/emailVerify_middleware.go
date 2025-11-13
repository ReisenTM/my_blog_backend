package middleware

import (
	"blog/internal/common/resp"
	"blog/internal/utils/email_store"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

type EmailVerifyMiddlewareRequest struct {
	Email     string `json:"email" binding:"required"`
	EmailCode string `json:"emailCode" binding:"required"`
}

// EmailVerifyMiddleware 邮箱验证码验证
func EmailVerifyMiddleware(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		resp.FailWithMsg("获取请求体错误", c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) //重置读取位置
	var cr EmailVerifyMiddlewareRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		logrus.Errorf("邮箱验证失败 %s", err)
		resp.FailWithMsg("邮箱验证失败", c)
		c.Abort()
		return
	}
	info, ok := email_store.Verify(cr.Email, cr.EmailCode)
	if !ok {
		fmt.Println("req:", cr.Email, cr.EmailCode)
		resp.FailWithMsg("邮箱验证码校验失败", c)
		c.Abort()
		return
	}
	c.Set("email", info.Email)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) //重置读取位置

}
