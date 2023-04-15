package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writers   []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()

	for _, w := range hook.Writers {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	}

	_, err := os.Stat("./logs")
	if os.IsNotExist(err) {
		err := os.Mkdir("logs", 0644)

		if err != nil {
			panic(err)
		}
	}

	allLogFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writers:   []io.Writer{allLogFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})
	l.SetLevel(logrus.TraceLevel)

	fmt.Println("Logger init success")
	e = logrus.NewEntry(l)
}
