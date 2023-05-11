package xlog

import (
	"WxProject/config"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = NewLogger()
}

func NewLogger() *logrus.Logger {
	Logger := logrus.New()
	// 日志级别
	logLev, err := logrus.ParseLevel(config.GetSystemConf().LogConf.LogLevel)
	Logger.SetLevel(logLev)
	if err != nil {
		Logger.SetLevel(logrus.DebugLevel)
	}
	// writer
	switch config.GetSystemConf().LogConf.LogOutPutMode {
	case "console":
		Logger.SetOutput(os.Stdout)
	case "file":
		path := config.GetSystemConf().LogConf.LogOutPutPath + "_%Y%m%d%H%M" + ""
		writer, _ := rotatelogs.New(
			path,
			rotatelogs.WithLinkName(config.GetSystemConf().LogConf.LogOutPutPath),
			rotatelogs.WithRotationCount(config.GetSystemConf().LogConf.LogFileRotationCount),
			rotatelogs.WithRotationTime(time.Hour*config.GetSystemConf().LogConf.LogRotationTime),
		)
		Logger.SetOutput(writer)
	case "both":
		{
			path := config.GetSystemConf().LogConf.LogOutPutPath + "_%Y%m%d%H%M" + ""
			writer, _ := rotatelogs.New(
				path,
				rotatelogs.WithLinkName(config.GetSystemConf().LogConf.LogOutPutPath),
				rotatelogs.WithRotationCount(config.GetSystemConf().LogConf.LogFileRotationCount),
				rotatelogs.WithRotationTime(time.Hour*config.GetSystemConf().LogConf.LogRotationTime),
			)
			multiWriter := io.MultiWriter(os.Stdout, writer)
			Logger.SetOutput(multiWriter)
		}
	default:
	}
	Logger.SetReportCaller(true)
	// 格式化
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceQuote:      true,
		TimestampFormat: config.GetSystemConf().LogConf.LogFileDateFmt,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})
	if config.GetSystemConf().LogConf.LogFormatter == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint:     true,
			TimestampFormat: config.GetSystemConf().LogConf.LogFileDateFmt,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := path.Base(frame.File)
				return frame.Function, fileName
			},
		})
	}
	return Logger
}
