package model

import (
	"blog/internal/model/enum"
	"gorm.io/gorm"
)

type LogModel struct {
	gorm.Model
	Type        enum.LogType   `json:"type"`  //日志类型
	Level       enum.LogLevel  `json:"level"` //日志级别
	Title       string         `gorm:"size:64" json:"title"`
	Content     string         `json:"content"`
	UserID      uint           `json:"user_id"` //可能是访客，id设为0
	UserModel   UserModel      `gorm:"foreignKey:UserID" json:"-"`
	IP          string         `gorm:"size:32" json:"ip"`
	Location    string         `gorm:"size:64" json:"addr"`
	IsRead      bool           `json:"is_read"`                     //是否已读
	LoginStatus bool           `json:"login_status"`                //登录状态
	Username    string         `gorm:"size:64" json:"username"`     //登录日志用户名
	Password    string         `gorm:"size:64" json:"password"`     //登录密码
	LoginType   enum.LoginType `gorm:"size:8" json:"login_type"`    //登录类型
	ServiceName string         `gorm:"size:32" json:"service_name"` //服务类型
}
