package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Create a global logger instance
var Log *logrus.Logger

func Init(filename string, level logrus.Level) {
	// Initialize the logger
	Log = logrus.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatal(err)
	}

	Log.SetOutput(file)


	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.SetLevel(level)
}
