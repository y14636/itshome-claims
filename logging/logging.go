package logging

import (
	. "os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitializeLogging(logFile string) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	log.SetFormatter(&log.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
}
