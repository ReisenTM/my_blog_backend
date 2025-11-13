package email_service

import (
	"blog/internal/global"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"strings"
)

const passTime = 10

// SendRegCode 发送注册验证码
func SendRegCode(to string, code string) error {
	subject := fmt.Sprintf("[%s]账号注册", global.Config.Email.SendNickname)
	content := fmt.Sprintf("你正在进行账号注册，这是你的验证码:%s,%d分钟内有效", code, passTime)
	return sendEmail(to, subject, content)
}

// SendResetCode 发送重置密码验证码
func SendResetCode(to string, code string) error {
	subject := fmt.Sprintf("[%s]密码找回", global.Config.Email.SendNickname)
	content := fmt.Sprintf("你正在进行密码重置，这是你的验证码:%s,%d分钟内有效", code, passTime)
	return sendEmail(to, subject, content)
}

// 发送邮件
func sendEmail(to string, subject string, content string) error {
	e := email.NewEmail()
	conf := global.Config.Email
	//设置发送方的邮箱
	e.From = fmt.Sprintf("%s <%s>", conf.SendNickname, conf.SendEmail)
	// 设置接收方的邮箱
	e.To = []string{to}
	//设置主题
	e.Subject = subject
	//设置文件发送的内容
	e.Text = []byte(content)
	//设置服务器相关的配置
	err := e.Send(fmt.Sprintf("%s:%d", conf.Domain, conf.Port), smtp.PlainAuth("", global.Config.Email.SendEmail, "swckpzzxifwjdeic", "smtp.qq.com"))
	if err != nil && !strings.Contains(err.Error(), "short response") {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}
