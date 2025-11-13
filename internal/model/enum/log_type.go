package enum

type LogType int8

const (
	LogLoginType LogType = iota + 1
	LogActionTypes
	LogRuntimeType
)
