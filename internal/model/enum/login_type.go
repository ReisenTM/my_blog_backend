package enum

type LoginType int8

const (
	UserPwdType LoginType = iota + 1
	QQLoginType
	EmailLoginType
)
