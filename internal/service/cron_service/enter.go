package cron_service

import (
	"github.com/robfig/cron/v3"
	"time"
)

// Cron 库的作用类似于定时器 ，用来缓存-DB同步
func Cron() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	// 每天2点去同步文章数据
	// 每天2点去同步文章数据
	crontab.AddFunc("0 0 2 * * *", SyncArticle)
	crontab.AddFunc("0 30 2 * * *", SyncUser)
	crontab.AddFunc("0 0 3 * * *", SyncComment)
	crontab.AddFunc("0 59 23 * * *", SyncSiteFlow)

	crontab.Start()
}
