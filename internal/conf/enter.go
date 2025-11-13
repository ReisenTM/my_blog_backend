package conf

type Config struct {
	System System `yaml:"system"`
	Log    Log    `yaml:"log"`
	DB     []DB   `yaml:"db"`    //连接的数据库
	Jwt    Jwt    `yaml:"jwt"`   //JWT
	Redis  Redis  `yaml:"redis"` //redis
	Email  Email  `yaml:"email"`
}
