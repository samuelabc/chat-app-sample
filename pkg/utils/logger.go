package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Define log directory and file
	logDir := "log"
	logFile := filepath.Join(logDir, "chat-app.log")

	// Check if log directory exists, if not create it
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			Log.Warn("Failed to create log directory, using default stderr")
			return
		}
	}

	// Open or create the log file within the log directory
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Warn("Failed to log to file, using default stderr")
		return
	}

	Log.SetOutput(file)
	Log.SetLevel(logrus.InfoLevel)
}
