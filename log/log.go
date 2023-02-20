package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.Formatter = new(logrus.JSONFormatter)
	Logger.Formatter = new(logrus.TextFormatter)                     //default
	Logger.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	Logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	Logger.Level = logrus.TraceLevel
	Logger.Out = os.Stdout
}
