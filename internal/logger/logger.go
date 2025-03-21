package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	DPanicln(args ...interface{})
	Panicln(args ...interface{})
	Fatalln(args ...interface{})

	Sync() error
}

func New(opts ...zap.Option) Logger {
	cfg := zap.NewProductionConfig()

	logger, err := cfg.Build(opts...)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger.Sugar()
}

func Fake() Logger {
	return zap.NewNop().Sugar()
}
