package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zonewave/copyer/log/jlog"
)

func InitLogger() {
	logInitOpt := func(logger *logrus.Logger) {
		logger.SetOutput(os.Stdout)
	}

	logrus.SetFormatter(&logrus.TextFormatter{})
	logInitOpt(logrus.StandardLogger())
	// init json logger
	jlog.InitJsonLogger()

}
