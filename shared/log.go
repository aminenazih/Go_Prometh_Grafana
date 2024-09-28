package shared

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes the logger based on the log level from the config
func InitLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warn("Invalid log level, defaulting to Info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter to JSON
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
