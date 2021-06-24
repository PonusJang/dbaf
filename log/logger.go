package log

import (
	log "github.com/sirupsen/logrus"
	"os"
)

/*
 * 日志模块
 */

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func Info(msg ...interface{}) {
	log.Info(msg)
}

func Debug(msg ...interface{}) {
	log.Debug(msg)
}

func Warn(msg ...interface{}) {
	log.Warn(msg)
}

func Fatal(msg ...interface{}) {
	log.Fatal(msg)
}

func Trace(msg ...interface{}) {
	log.Trace(msg)
}

func Panic(msg ...interface{}) {
	log.Panic(msg)
}

func Error(msg ...interface{}) {
	log.Error(msg)
}
