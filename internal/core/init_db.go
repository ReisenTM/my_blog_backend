package core

import (
	"blog/internal/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitDb() *gorm.DB {
	if len(global.Config.DB) <= 0 {
		logrus.Errorf("未配置数据库")
		return nil
	}
	dc := global.Config.DB[0] //0库主库：写 ，其余从库 读
	if dc.Host == "" {
		logrus.Warnln("未配置数据库连接地址")
	}
	dsn := dc.Dsn()
	var myLogger logger.Interface
	//gorm连接
	if dc.Debug == true {
		//Debug环境下显示log
		myLogger = logger.Default.LogMode(logger.Info)
	} else {
		//正常模式下仅显示错误
		myLogger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
		Logger:                                   myLogger,
	})
	if err != nil {
		logrus.Fatalf("连接数据库失败 %s", err)
	}
	//拿到原始sqlDB
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)  //最大空闲连接数
	sqlDB.SetMaxOpenConns(100) //最多可容纳
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infoln("数据库 连接成功")

	return db
}
