package rockgo

import (
	"log"
	"os"
)

var Debug = true

var rockLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func Debug(v ...interface{}) {
	rockLogger.Println(v)
}

func Info(v ...interface{}) {
	rockLogger.Println(v)
}

func Warn(v ...interface{}) {
	rockLogger.Println(v)
}

func Error(v ...interface{}) {
	rockLogger.Println(v)
}
