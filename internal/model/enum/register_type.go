package enum

type RegisterSourceType int8

// 注册来源
const (
	RegEmailSourceType    RegisterSourceType = 1
	RegQQSourceType       RegisterSourceType = 2
	RegGitSourceType      RegisterSourceType = 3
	RegTerminalSourceType RegisterSourceType = 4 //命令行创建
)
