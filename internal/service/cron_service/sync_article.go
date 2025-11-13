package cron_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/redis_service/redis_article"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SyncArticle 文章定时同步
func SyncArticle() {
	collectMap := redis_article.GetAllCacheCollect()
	favorMap := redis_article.GetAllCacheFavor()
	viewMap := redis_article.GetAllCacheLook()
	commentMap := redis_article.GetAllCacheComment()

	var list []model.ArticleModel
	global.DB.Find(&list)

	for _, article := range list {
		collect := collectMap[article.ID]
		favor := favorMap[article.ID]
		view := viewMap[article.ID]
		comment := commentMap[article.ID]
		if collect == 0 || favor == 0 || view == 0 || comment == 0 {
			continue
		}

		err := global.DB.Model(&article).Updates(map[string]any{
			"views_count":   gorm.Expr("views_count + ?", view),
			"favor_count":   gorm.Expr("favor_count + ?", favor),
			"collect_count": gorm.Expr("collect_count + ?", collect),
			"comment_count": gorm.Expr("comment_count + ?", comment),
		}).Error
		if err != nil {
			logrus.Errorf("更新失败 %s", err)
			continue
		}
		logrus.Infof("%s 更新成功", article.Title)
	}

	// 这里可能会有增量的数据
	// 可以再获取一次
	//collectMap := redis_article.GetAllCacheCollect()
	//favorMap := redis_article.GetAllCacheDigg()
	//lookMap := redis_article.GetAllCacheLook()

	// 走完之后清空掉缓存
	redis_article.Clear()

	// 再同步回去

}
