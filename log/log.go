package log

import (
	"fmt"
	"github.com/robfig/cron"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	lock     sync.Mutex
	filename string // should be set to the actual filename
	fp       *os.File
}
var timeFormat = "2006-01-02"
var logger *Logger
var c = cron.New()

// Make a new RotateWriter. Return nil if error occurs during setup.
func GetLogger(Mode string) *Logger {
	if logger != nil {
		return logger
	}
	logger = newLogger(Mode)
	return logger
}

func newLogger(Mode string) *Logger {
	if Mode == "prod"{
		logger = &Logger{filename:"stdout"}
		fmt.Println("PROD")
		logger.Run()
		c.AddFunc("* * * * *", logger.Run)
		c.Start()
	}else{
		logger = &Logger{filename:"stdout",fp:os.Stdout}
	}
	return logger
}

func (l *Logger)Info(Msg string){
	fmt.Println("FILE NAME",l.filename)
	log := []string{"INFO",Msg,time.Now().Format(time.RFC3339) + "\n"}
	l.Write(strings.Join(log, "|"))
}

func (l *Logger)Error(Msg string){
	log := []string{"ERROR",Msg,time.Now().Format(time.RFC3339) + "\n"}
	l.Write(strings.Join(log, "|"))
}
// Write satisfies the io.Writer interface.
func (l *Logger) Write(msg string) (int, error) {

	l.lock.Lock()
	defer l.lock.Unlock()
	return l.fp.Write([]byte(msg))
}

// Perform the actual act of rotating and reopening file.
func (l *Logger) Run() {
	l.lock.Lock()
	defer l.lock.Unlock()
	// 파일 닫고.
	if l.fp != nil {
		err := l.fp.Close()
		l.fp = nil
		if err != nil {
			return
		}
	}

	// 파일이름 지정
	l.filename = time.Now().Format(timeFormat) + ".log"
	// 로그파일이 이미 있으면.
	_, err := os.Stat(l.filename)
	if err == nil {
		l.fp, err = os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY,7777)
		return
	}
	// 새로운 로그파일 생성.
	// Create a file.
	l.fp, err = os.Create(l.filename)
	return
}