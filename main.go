package main

import (
	"blog/internal/core"
	"blog/internal/flags"
	"blog/internal/global"
	"blog/internal/routers"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	flags.Parse()                   //绑定命令行参数
	global.Config = core.ReadConf() //读配置文件
	core.InitDefaultLogus()         //初始化日志
	global.DB = core.InitDb()       //初始化数据库
	global.Redis = core.InitRedis() //初始化redis
	flags.Run()
	//启动程序
	routers.Run()
}
