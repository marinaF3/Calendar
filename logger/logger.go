package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
)

func Init(filename string) error {
	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func Close() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			return
		}
	}
}

func Info(msg string) {
	err := infoLogger.Output(2, msg)
	if err != nil {
		return
	}
}

func Error(msg string) {
	err := errorLogger.Output(2, msg)
	if err != nil {
		return
	}
}
