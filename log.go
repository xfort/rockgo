package rockgo

import (
	"log"
	"os"
	"fmt"
	"errors"
)

var DebugTag bool = true

var rockLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

type RockLog struct {
	log.Logger
}

func (rocklog *RockLog) Warn(v ...interface{}) {
	rocklog.Println(v)
	//log.Println()
}

func Debug(v ...interface{}) {
	rockLogger.Output(2, fmt.Sprintln("debug", v))
}

func Info(v ...interface{}) {
	rockLogger.Output(2, fmt.Sprintln("info", v))
}

func Warn(v ...interface{}) {
	rockLogger.Output(2, fmt.Sprintln("warn", v))

}

func Error(v ...interface{}) {
	rockLogger.Output(2, fmt.Sprintln("error", v))
}

func NewError(v ...interface{}) error {
	return errors.New(fmt.Sprint(v))
}
