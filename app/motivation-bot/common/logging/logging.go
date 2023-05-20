package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()

	for _, writer := range hook.Writer {
		writer.Write([]byte(line))
	}

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry // create pointer to logrus

type Logger struct { // create a class with logger
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLoggerWithFields(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

func init() { // @notice if init with small "i" this is function will be run automatically when somewhere will be import this module
	l := logrus.New()
	l.SetReportCaller(true)

	formatter := &logrus.TextFormatter{CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		fileName := path.Base(frame.File)

		return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", fileName, frame.Line)
	}, FullTimestamp: true, ForceColors: true}

	l.SetFormatter(formatter)

	//err := os.MkdirAll("logs", 0777)
	//common.CriticalErrorHandler(err)
	//
	//file, err := os.OpenFile("logs/logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	//common.CriticalErrorHandler(err)

	//defer file.Close()

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		//Writer:    []io.Writer{file, os.Stdout},
		Writer:    []io.Writer{os.Stdout},
		LogLevels: logrus.AllLevels,
	})
	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l) // set to our E variable from above scope
}
