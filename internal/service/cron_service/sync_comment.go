package cron_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/redis_service/redis_comment"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SyncComment 评论定时同步
func SyncComment() {
	commentFavorMap := redis_comment.GetAllCacheFavor()
	commentReplyMap := redis_comment.GetAllCacheReply()
	var list []model.CommentModel
	global.DB.Find(&list)

	for _, comment := range list {
		favor := commentFavorMap[comment.ID]
		reply := commentReplyMap[comment.ID]
		if reply == 0 || favor == 0 {
			continue
		}

		err := global.DB.Model(&comment).Updates(map[string]any{
			"favor_count": gorm.Expr("favor_count + ?", favor),
			"reply_count": gorm.Expr("reply_count + ?", reply),
		}).Error
		if err != nil {
			logrus.Errorf("更新失败 %s", err)
			continue
		}
		logrus.Infof("%s 更新成功", comment.ID)
	}

	// 走完之后清空掉缓存
	redis_comment.Clear()
}
