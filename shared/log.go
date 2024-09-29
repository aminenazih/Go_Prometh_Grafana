package shared

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warn("Invalid log level, defaulting to Info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
