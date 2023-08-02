package jlog

import (
	"github.com/sirupsen/logrus"
)

var JLog *logrus.Logger

func InitJsonLogger(opts ...func(*logrus.Logger)) {
	JLog = NewJsonLogger(opts...)
}

func NewJsonLogger(opts ...func(*logrus.Logger)) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}
