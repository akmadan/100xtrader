package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

// InitLogger configures logrus for production use
func InitLogger() {
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
		PrettyPrint:     false,
	})
	Logger.SetOutput(os.Stdout)

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}
	level, err := logrus.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
}
