package log

import (
	"fmt"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/evalphobia/logrus_sentry"
)

const (
	InfoLevelLog  = log.InfoLevel
	ErrorLevelLog = log.ErrorLevel
	FatalLevelLog = log.FatalLevel
	PanicLevelLog = log.PanicLevel
)

var logger *log.Logger
var once sync.Once
var topic = "go-core"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func getLoggerInstance(dsn string, sentry bool) *log.Logger {
	once.Do(func() {
		logger = log.New()

		if sentry {
			// implement hook for sentry
			if hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
				log.PanicLevel,
				log.FatalLevel,
				log.ErrorLevel,
			}); err == nil {
				logger.Hooks.Add(hook)
			}
		}
	})

	return logger
}

func logContext(topic string) *log.Entry {
	return logger.WithFields(log.Fields{
		"topic": topic,
		"path":  getCtx(),
	})
}

func Init(dsn string, sentry bool) {
	logger = getLoggerInstance(dsn, sentry)
}

func New(level log.Level, topic string, message ...interface{}) {
	entry := logContext(topic)
	switch level {
	case log.DebugLevel:
		entry.Debug(message...)
	case log.InfoLevel:
		entry.Info(message...)
	case log.WarnLevel:
		entry.Warn(message...)
	case log.ErrorLevel:
		entry.Error(message...)
	case log.FatalLevel:
		entry.Fatal(message...)
	case log.PanicLevel:
		entry.Panic(message...)
	}
}

func getCtx() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[2])
	file, line := f.FileLine(pc[2])
	return fmt.Sprintf("%s:%d", file, line)
}
