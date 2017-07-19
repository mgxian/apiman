package log

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/pkg/setting"
)

func LoggerInit() {
	sec := setting.Cfg.Section("log")

	logFile := sec.Key("file").String()
	logLevel := sec.Key("level").String()

	// log to file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		// set log level
		switch logLevel {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		default:
			log.Warning("unsport log level user default info")
			log.SetLevel(log.InfoLevel)
		}

		// set log file
		log.SetOutput(file)
	} else {
		log.Warning("Failed to log to file, using default stderr")
	}
}
