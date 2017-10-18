package pvhelpers

import (
	"fmt"
	"log"
)

func LogVerbose(format string, v ...interface{}) {
	if Cfg.App.LogLevel == "Verbose" {
		log.Printf(format, v...)
	}
}

func LogInfo(format string, v ...interface{}) {
	if Cfg.App.LogLevel == "Verbose" || Cfg.App.LogLevel == "Info" {
		log.Printf(format, v...)
	}
}

func LogWarn(format string, v ...interface{}) {
	if Cfg.App.LogLevel == "Verbose" || Cfg.App.LogLevel == "Info" || Cfg.App.LogLevel == "Warn" {
		log.Printf(format, v...)
	}
}

func LogError(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func LogErrorObject(err error) {
	LogError(fmt.Sprintf("Error: ", err))
}
