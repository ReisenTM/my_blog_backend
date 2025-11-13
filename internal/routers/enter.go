package routers

import (
	"blog/internal/global"
	"blog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run() {
	//设置运行模式
	gin.SetMode(global.Config.System.GinMode)

	r := gin.Default()
	//路径请求映射
	r.Static("/uploads", "uploads")
	//创建路由组
	gr := r.Group("/api")
	//使用全局中间件
	gr.Use(middleware.LogMiddleWare).Use(middleware.Cors())
	UserRouters(gr)
	//启动路由监听
	addr := global.Config.System.Addr()
	err := r.Run(addr)
	if err != nil {
		logrus.Errorf("server启动失败:%v", err)
		return
	}
}
