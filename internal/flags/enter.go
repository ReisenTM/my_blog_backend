package flags

import (
	"flag"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
	Type    string //命令类型
	Sub     string //子命令
	ES      bool
}

var FlagOptions = new(Options)

// Parse flag绑定
func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FlagOptions.Version, "v", false, "版本")
	flag.StringVar(&FlagOptions.Type, "t", "", "类型")
	flag.StringVar(&FlagOptions.Sub, "s", "", "子类")
	flag.Parse()
}

// Run flag实现
func Run() {
	if FlagOptions.DB {
		//执行数据库迁移
		FlagDB()
		os.Exit(0)
	}

	switch FlagOptions.Type {
	case "user":
		u := FlagUser{}
		switch FlagOptions.Sub {
		case "create":
			//创建用户
			u.Create()
			os.Exit(0)
		}
	}
}
