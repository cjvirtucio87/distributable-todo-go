package rlogging

import "go.uber.org/zap"

type Logger interface {
	Infof(tmpl string, args ...interface{})
	Debugf(tmpl string, args ...interface{})
	Errorf(tmpl string, args ...interface{})
}

func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()

	defer logger.Sync()

	return &ZapLogger{
		SugaredLogger: logger.Sugar(),
	}
}
