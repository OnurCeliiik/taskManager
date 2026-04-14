package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	Log  *logrus.Logger
	once sync.Once
)

func init() {
	// Auto-initialize on first use
	ensureInit()
}

func ensureInit() {
	once.Do(func() {
		Log = logrus.New()

		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		Log.SetOutput(os.Stdout)
		Log.SetLevel(logrus.InfoLevel)

		if os.Getenv("LOG_LEVEL") == "debug" {
			Log.SetLevel(logrus.DebugLevel)
		}
	})
}

func Init() {
	ensureInit()
}

func Info(msg string, fields ...interface{}) {
	ensureInit()
	Log.WithFields(parseFields(fields...)).Info(msg)
}

func Error(msg string, fields ...interface{}) {
	ensureInit()
	Log.WithFields(parseFields(fields...)).Error(msg)
}

func Debug(msg string, fields ...interface{}) {
	ensureInit()
	Log.WithFields(parseFields(fields...)).Debug(msg)
}

func Warn(msg string, fields ...interface{}) {
	ensureInit()
	Log.WithFields(parseFields(fields...)).Warn(msg)
}

func parseFields(fields ...interface{}) logrus.Fields {
	customFields := logrus.Fields{}
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			customFields[fields[i].(string)] = fields[i+1]
		}
	}
	return customFields
}
