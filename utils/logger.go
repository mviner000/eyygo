package utils

import (
    "log"
    "os"
)

func GetLogger(logFile string) *log.Logger {
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }
    return log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}