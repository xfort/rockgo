package rockgo

import (
	"log"
	"os"
	"path/filepath"
	"errors"
	"github.com/xfort/rockgo/proto"
	"time"
	"fmt"
	"encoding/json"
)

type RockLogIn interface {
	Log(obj *proto.LogObj)
	Debug(tag string, v ...interface{})
	Info(tag string, v ...interface{})
	Warn(tag string, v ...interface{})
	Error(tag string, v ...interface{})
	Fatal(tag string, v ...interface{})
}

type RockLogger struct {
	goLogger *log.Logger
	logFile  string
}

var Defaultlogger *RockLogger

func init() {
	Defaultlogger = NewRockLogger("rock")
}

func NewRockLogger(logfilepre string) *RockLogger {
	dirFile, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dirFile = os.TempDir()
	}
	dirFile = filepath.Join(dirFile, "log")

	rockLogger := &RockLogger{}
	filename := time.Now().Format("2006-01-02")
	logFILE := filepath.Join(dirFile, logfilepre+filename+".log")
	rockLogger.InitData(logFILE)
	return rockLogger
}

func (rl *RockLogger) InitData(logfile string) error {
	rl.logFile = logfile

	err := os.MkdirAll(filepath.Dir(logfile), 0666)
	if err != nil {
		return errors.New("创建logfile文件夹失败" + err.Error() + logfile)
	}
	logFileObj, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("打开logfile失败" + logfile + err.Error())
	}

	rl.goLogger = log.New(logFileObj, "", log.Ldate|log.Ltime|log.Llongfile)
	return nil
}

func (rl *RockLogger) Log(logObj *proto.LogObj) {
	jsonByte, err := json.Marshal(logObj)
	if err != nil {
		log.Println("LogObj转为json失败", err.Error())
	}
	rl.goLogger.Output(3, fmt.Sprintln(string(jsonByte))+fmt.Sprintln())
}

func (rl *RockLogger) LogMsg(lv proto.LogLevel, tag string, msg string) {
	logObj := &proto.LogObj{Level: lv, Tag: tag, Message: msg, TimestampUTC: time.Now().UTC().Unix()}
	rl.Log(logObj)
}

func (rl *RockLogger) Debug(tag string, v ...interface{}) {
	logObj := &proto.LogObj{Level: proto.LogLevel_Debug, Tag: tag, Message: fmt.Sprint(v...), TimestampUTC: time.Now().UTC().Unix()}
	rl.Log(logObj)
	log.Println(tag, v)
}

func (rl *RockLogger) Info(tag string, v ...interface{}) {
	logObj := &proto.LogObj{Level: proto.LogLevel_Info, Tag: tag, Message: fmt.Sprint(v...), TimestampUTC: time.Now().UTC().Unix()}
	rl.Log(logObj)
	log.Println(tag, v)
}
func (rl *RockLogger) Warn(tag string, v ...interface{}) {
	logObj := &proto.LogObj{Level: proto.LogLevel_Warn, Tag: tag, Message: fmt.Sprint(v...), TimestampUTC: time.Now().UTC().Unix()}
	rl.Log(logObj)
	log.Println(tag, v)
}
func (rl *RockLogger) Error(tag string, v ...interface{}) {
	logObj := &proto.LogObj{Level: proto.LogLevel_Error, Tag: tag, Message: fmt.Sprint(v...), TimestampUTC: time.Now().UTC().Unix()}
	rl.Log(logObj)
	log.Println(tag, v)
}
func (rl *RockLogger) Fatal(tag string, v ...interface{}) {
	logObj := &proto.LogObj{Level: proto.LogLevel_Fatal, Tag: tag, Message: fmt.Sprint(v...), TimestampUTC: time.Now().UTC().Unix()}
	rl.Log(logObj)
	log.Println(tag, v)
}
func Debug(tga string, v ...interface{}) {
	Defaultlogger.LogMsg(proto.LogLevel_Debug, tga, fmt.Sprint(v))
}

func Info(tga string, v ...interface{}) {
	Defaultlogger.LogMsg(proto.LogLevel_Info, tga, fmt.Sprint(v))
}

func Warn(tga string, v ...interface{}) {
	Defaultlogger.LogMsg(proto.LogLevel_Warn, tga, fmt.Sprint(v))
}

func Error(tga string, v ...interface{}) {
	Defaultlogger.LogMsg(proto.LogLevel_Error, tga, fmt.Sprint(v))
}

func Fatal(tga string, v ...interface{}) {
	Defaultlogger.LogMsg(proto.LogLevel_Fatal, tga, fmt.Sprint(v))
}
