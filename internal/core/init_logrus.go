package core

import (
	"blog/internal/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 前缀结构体
type PrefixFormatter struct {
	Prefix    string
	Formatter logrus.Formatter
}

// Format 重写接口，设置前缀
func (p *PrefixFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = entry.Message
	entry.Data["prefix"] = p.Prefix
	return p.Formatter.Format(entry)
}

type CustomFormatter struct {
	Prefix string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	levelColor := ""
	resetColor := "\033[0m"
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = "\033[34m" // 蓝色
	case logrus.InfoLevel:
		levelColor = "\033[32m" // 绿色
	case logrus.WarnLevel:
		levelColor = "\033[33m" // 黄色
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = "\033[31m" // 红色
	default:
		levelColor = ""
	}
	coloredLevel := fmt.Sprintf("%s%s%s", levelColor, level, resetColor)
	file := ""
	if entry.HasCaller() {
		file = fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
	}
	msg := fmt.Sprintf("%s %s [%s] %s %s\n", f.Prefix, timestamp, coloredLevel, file, entry.Message)
	return []byte(msg), nil
}

// InitDefaultLogus 初始化默认 logrus 配置，带时间、行号、颜色
func InitDefaultLogus() {
	// 创建 logs/2025-05-01 文件夹
	today := time.Now().Format("2006-01-02")
	logDir := filepath.Join("logs", today)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 构造日志文件路径
	logFile := filepath.Join(logDir, global.Config.Log.App+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("打开日志文件失败: %v", err)
	}
	// 设置 logrus 输出
	mw := io.MultiWriter(file, os.Stdout)
	logrus.SetOutput(mw)
	logrus.SetReportCaller(true) // 开启调用者信息（行号）
	// 自定义格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", frame.File + ":" + strconv.Itoa(frame.Line)
		},
	})

	logrus.SetLevel(logrus.DebugLevel) // 默认最低等级为 Debug
	logrus.Info("日志 初始化完成")
}
