package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func Init() {
	logOutput, err := os.OpenFile("bvpn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logOutput = os.Stdout
	}
	infoLogger = log.New(logOutput, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	errorOutput, err := os.OpenFile("bvpn-errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		errorOutput = os.Stderr
	}
	errorLogger = log.New(errorOutput, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Log(message string) {
	infoLogger.Println(message)
}

func LogError(message string) {
	errorLogger.Println(message)
}
