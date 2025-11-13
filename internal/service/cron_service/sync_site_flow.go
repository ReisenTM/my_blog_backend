package cron_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/redis_service/redis_site"
)

// SyncSiteFlow 同步网站访问量
func SyncSiteFlow() {
	flow := redis_site.GetFlow()
	global.DB.Create(&model.SiteFlowModel{Count: flow})
	redis_site.ClearFlow()
}
