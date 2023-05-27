package logr

import (
	"log"

	"go.uber.org/zap"
)

// Logger is an interface to print logs
type Logger interface {
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Errorf(format string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

var DefaultLogger Logger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	DefaultLogger = logger.Sugar()
}
