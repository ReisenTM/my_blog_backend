package conf

import "fmt"

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Debug    bool   `yaml:"debug"`  //是否打印全部日志
	Source   string `yaml:"source"` //数据库的源 pgsql mysql
}

// Dsn 拼接dsn 主库
// Mysql
//
//	func (d DB) Dsn() string {
//		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//			d.User, d.Password, d.Host, d.Port, d.Database)
//	}
func (d DB) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Database)
}

// Empty 判空
func (d *DB) Empty() bool {
	return d.User == "" && d.Password == "" && d.Host == "" && d.Database == "" && d.Source == ""
}

// Addr 地址
func (d DB) Addr() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}
