package rockgo

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type RockLogger struct {
	goLogger *log.Logger
	logFile  string
}

var Defaultlogger = NewRockLogger()

func NewRockLogger() *RockLogger {
	dirFile, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dirFile = os.TempDir()
	}
	dirFile = filepath.Join(dirFile, "log")

	rockLogger := &RockLogger{}
	filename := time.Now().Format("2006-01-02_15_04")
	logFILE := filepath.Join(dirFile, filename+strconv.FormatInt(time.Now().Unix(), 10)+".log")
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



