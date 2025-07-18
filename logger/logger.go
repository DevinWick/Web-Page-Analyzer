package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}
