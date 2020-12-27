package logging

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() (err error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if !gin.IsDebugging() {
		config.Level.SetLevel(zapcore.WarnLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return
	}

	zap.ReplaceGlobals(logger)

	return
}

func Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}

func Error(args ...interface{}) {
	zap.S().Error(args...)
}
