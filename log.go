package rockgo

import (
	"log"
	"os"
	"fmt"
	"errors"
	"time"
	"sync"
)

const (
	Log_Debug = 1
	Log_Info  = 2
	Log_Warn  = 4
	Log_Error = 8
	Log_Fatal = 16
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

var logMsgPool sync.Pool = sync.Pool{New: func() interface{} {
	return &LogMsg{}
}}

type LogMsg struct {
	Id  int64
	LV  int
	Tag string
	Msg string
}

func ObtainLogMsg(tag string, lv int, v ...interface{}) *LogMsg {
	logmsg := logMsgPool.Get().(*LogMsg)
	logmsg.Id = time.Now().Unix()
	logmsg.LV = lv
	logmsg.Tag = tag
	logmsg.Msg = fmt.Sprint(v)
	return logmsg
}

func recycleLoMsg(logmsg *LogMsg) {
	logmsg.Id = 0
	logmsg.LV = 0
	logmsg.Tag = ""
	logmsg.Msg = ""
	logMsgPool.Put(logmsg)
}
