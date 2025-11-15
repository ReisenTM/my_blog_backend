package flags

import (
	"blog/internal/global"
	"blog/internal/model"
	"github.com/sirupsen/logrus"
)

// FlagDB 迁移数据库
func FlagDB() {

	err := global.DB.AutoMigrate(
		&model.UserModel{},
		&model.LogModel{},
		&model.ArticleModel{})
	if err != nil {
		logrus.Errorf("自动迁移失败 %s", err)
		return
	}
	logrus.Infof("自动迁移成功")
}
