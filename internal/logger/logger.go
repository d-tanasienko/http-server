package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Logger struct{}

func (logger Logger) Error(message string) {
	log.Println("ERROR: " + message)
}

func NewLogger() *Logger {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	logsDir := filepath.Join(basepath, "..", "..", "logs", "error.log")
	if _, err := os.Stat(logsDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(logsDir, 0755)
		} else {
			log.Fatal(err)
		}
	}
	file, err := os.OpenFile(filepath.Join(logsDir, "error.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	return &Logger{}
}
