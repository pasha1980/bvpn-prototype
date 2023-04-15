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
	output, err := os.OpenFile("bvpn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		output = os.Stdout
	}

	infoLogger = log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Log(message string) {
	infoLogger.Println(message)
}

func LogError(message string) {
	errorLogger.Println(message)
}
