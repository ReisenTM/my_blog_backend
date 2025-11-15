package model

import (
	"blog/internal/model/enum"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username  string                  `gorm:"unique;not null" json:"username"`
	Email     string                  `gorm:"unique" json:"email"`
	Password  string                  `gorm:"not null" json:"password"`
	Avatar    string                  `gorm:"not null" json:"avatar"`
	Status    int                     `gorm:"default:1" json:"status"` //1 正常 2 Baned
	RegSource enum.RegisterSourceType `json:"regSource"`               //注册来源
	Role      enum.RoleType           `json:"role"`
}
