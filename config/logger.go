package config

import (
	"io"
	"log"
	"os"
)

var debugLogger *log.Logger

func init() {
	// Initialize debugLogger immediately with default output
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func InitLogger(debug bool) {
	if !debug {
		debugLogger.SetOutput(io.Discard)
		log.SetOutput(io.Discard)
	}
}

func DebugLog(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}
