package global

import (
	"blog/internal/conf"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// global用来为项目内的方法提供操作对象
var (
	Config *conf.Config  //全局配置
	DB     *gorm.DB      //数据库
	Redis  *redis.Client //redis
)
