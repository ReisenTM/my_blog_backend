package enum

type LoginType int8

const (
	EmailLoginType LoginType = iota + 1
	QQLoginType
	GithubLoginType
)
