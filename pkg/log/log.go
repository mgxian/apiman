package log

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/pkg/setting"
)

var (
	Log *logrus.Logger
)

func GetLogger() (*logrus.Logger, error) {
	if Log != nil {
		return Log, nil
	}
	return nil, errors.New("Log is nil")
}

//func LoggerInit() (*logrus.Logger, error) {
func LoggerInit() {
	// create logger
	Log = logrus.New()

	//get log config
	//setting.NewConfig()
	sec := setting.Cfg.Section("log")

	logFile := sec.Key("file").String()
	logLevel := sec.Key("level").String()

	// log to file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		// set log level
		switch logLevel {
		case "debug":
			Log.Level = logrus.DebugLevel
		case "info":
			Log.Level = logrus.InfoLevel
		case "warn":
			Log.Level = logrus.WarnLevel
		case "error":
			Log.Level = logrus.ErrorLevel
		default:
			Log.Info("unsport log level user default info")
			Log.Level = logrus.InfoLevel
		}

		// set log file
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

}
