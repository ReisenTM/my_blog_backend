package enum

type LogLevel int8

const (
	LogLevelInfo LogLevel = iota + 1
	LogLevelWarn
	LogLevelError
)

func (level LogLevel) String() string {
	switch level {
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	}
	return ""
}
