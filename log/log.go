package log

import (
	"fmt"
	"runtime"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_sentry"
	config "github.com/spf13/viper"
)

const (
	InfoLevelLog  = log.InfoLevel
	ErrorLevelLog = log.ErrorLevel
	FatalLevelLog = log.FatalLevel
	PanicLevelLog = log.PanicLevel
)

var instance *log.Logger
var once sync.Once

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func getLoggerInstance() *log.Logger {
	once.Do(func() {
		instance = log.New()

		// implement hook for sentry
		dsn := config.GetString("sentry.dsn")
		hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})

		if err == nil {
			instance.Hooks.Add(hook)
		}
	})

	return instance
}

func logContext() *log.Entry {
	if config.GetBool("sentry.enabled") {
		logger := getLoggerInstance()
		return logger.WithFields(log.Fields{
			"topic": "digital-products",
			"path":  getCtx(),
		})
	}

	return log.WithFields(log.Fields{
		"topic": "digital-products",
		"path":  getCtx(),
	})
}

func New(level log.Level, message ...interface{}) {
	entry := logContext()
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
