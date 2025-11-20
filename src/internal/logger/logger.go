package logger

import (
	"log"
)

// Simple wrapper — swap to more advanced logger later
func Info(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

func Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func Debug(format string, v ...interface{}) {
	log.Printf("[DEBUG] "+format, v...)
}
