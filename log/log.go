package log

import (
	"fmt"

	"go.uber.org/zap"
)

var DefaultLogger *zap.Logger

func init() {
	var err error
	DefaultLogger, err = zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
}

func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Debug(fmt.Sprintf(msg, args...))
}
func Debug(msg string, fields ...zap.Field) {
	DefaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	DefaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	DefaultLogger.Warn(msg, fields...)
}

func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Error(fmt.Sprintf(msg, args...))
}
func Error(msg string, fields ...zap.Field) {
	DefaultLogger.Error(msg, fields...)
}
