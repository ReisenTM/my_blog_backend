package core

import (
	"blog/internal/conf"
	"blog/internal/flags"
	"blog/internal/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

// ReadConf 读取配置文件
func ReadConf() *conf.Config {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	c := new(conf.Config)
	err = yaml.Unmarshal(byteData, c)
	if err != nil {
		panic(fmt.Errorf("yaml配置文件格式错误%v", err))
	}
	fmt.Printf("读取配置文件成功:%s\n", flags.FlagOptions.File)
	return c
}

// SaveConf 配置文件更新
func SaveConf() {
	c := global.Config
	data, err := yaml.Marshal(c)
	if err != nil {
		logrus.Errorf("配置文件转换失败,err:%v", err)
		return
	}
	err = os.WriteFile(flags.FlagOptions.File, data, 0666)
	if err != nil {
		logrus.Errorf("配置文件写入失败,err:%v", err)
		return
	}
	logrus.Info("配置文件更新成功")
}
