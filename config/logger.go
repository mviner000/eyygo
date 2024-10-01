package config

import (
	"io"
	"log"
	"os"
)

var debugLogger *log.Logger

func init() {
	// Initialize debugLogger immediately
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func initLogger() {
	if !AppSettings.Debug {
		debugLogger.SetOutput(io.Discard)
		log.SetOutput(io.Discard)
	}
}

func DebugLog(format string, v ...interface{}) {
	if AppSettings.Debug {
		debugLogger.Printf(format, v...)
	}
}
