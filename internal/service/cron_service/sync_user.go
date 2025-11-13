package cron_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/redis_service/redis_user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncUser() {
	lookMap := redis_user.GetAllCacheLook()

	var list []model.UserConfModel
	global.DB.Find(&list)

	for _, m := range list {
		look := lookMap[m.UserID]
		if look == 0 {
			continue
		}

		err := global.DB.Model(&m).Updates(map[string]any{
			"views_count": gorm.Expr("views_count + ?", look),
		}).Error
		if err != nil {
			logrus.Errorf("更新失败 %s", err)
			continue
		}
		logrus.Infof("%s 更新成功", m.UserID)
	}

	// 走完之后清空掉
	redis_user.Clear()

	// 再同步回去
}
