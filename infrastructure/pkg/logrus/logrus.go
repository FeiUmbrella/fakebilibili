package logrus

import (
	"encoding/json"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

// JsonInfo 输出日志的JSON结构体
type JsonInfo struct {
	Time     string `json:"time"`
	Level    string `json:"level"`
	Msg      string `json:"msg"`
	File     string `json:"file,omitempty"` // omitempty-当该字段为0值时，在JSON中忽略该字段
	Function string `json:"function,omitempty"`
}

// JsonFormatter 自定义json解析
type JsonFormatter struct {
	logrus.JSONFormatter
}

// Format 接口实现，用于自定义日志格式
func (f *JsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	info := &JsonInfo{
		Time:  entry.Time.Format("2006-01-02 15:04:05"), // 日志时间戳格式
		Level: entry.Level.String(),                     // 日志级别
		Msg:   entry.Message,                            // 日志信息
	}

	// 日志级别匹配时打印调用者信息
	if entry.Level == logrus.DebugLevel || entry.Level == logrus.WarnLevel || entry.Level == logrus.ErrorLevel || entry.Level == logrus.PanicLevel {
		info.Function = entry.Caller.Function // 生成日志对应的函数位置
		info.File = entry.Caller.File         // 生成日志的对应文件
	}
	jsonData, err := json.Marshal(info) // 将结构体变为对应的JSON的字节流
	if err != nil {
		return nil, err
	}
	finalLog := string(jsonData) + "\n" // 为日志加上一个换行符
	return []byte(finalLog), nil
}

var Logger *logrus.Logger

func init() {
	Logger = ReturnInstance()
}

var logFilePath = "./runtime/log" // 日志文件保存路径

func ReturnInstance() *logrus.Logger {
	logger := logrus.New()
	// 日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 打印调用者信息
	logger.SetReportCaller(true)
	// 定义到空输出
	logger.SetOutput(io.Discard)

	// 设置 rotatelogs 实现文件分割设置
	logInfoWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d/info.log",         // 将info级别日志输出至对应文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    // 文件保存最长时间
		rotatelogs.WithRotationTime(1*time.Hour), // 一小时分割一次日志文件
	)
	logFatalWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d/fatal.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	logDebugWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d/debug.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	logWarnWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d/warn.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	logErrorWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d/error.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	logPanicWriter, _ := rotatelogs.New(
		logFilePath+"/%Y-%m-%d/panic.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)

	// 设置hook
	writerMap := lfshook.WriterMap{
		// 将不同级别日志实现对应的分割设置
		logrus.InfoLevel:  logInfoWriter,
		logrus.FatalLevel: logFatalWriter,
		logrus.DebugLevel: logDebugWriter,
		logrus.WarnLevel:  logWarnWriter,
		logrus.ErrorLevel: logErrorWriter,
		logrus.PanicLevel: logPanicWriter,
	}
	logger.Formatter = &JsonFormatter{} // 将日志以JSON类型输出，与之对应的以Text类型输出
	// 添加hook
	logger.AddHook(lfshook.NewHook(writerMap, &JsonFormatter{}))

	return logger
}
