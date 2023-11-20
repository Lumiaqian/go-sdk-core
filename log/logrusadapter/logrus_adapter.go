package logrusadapter

import (
	"context"

	"github.com/Lumiaqian/go-sdk-core/log"

	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
	logrusLogger *logrus.Logger
}

func NewLogrusAdapter(logrusLogger *logrus.Logger) log.Logger {
	return &LogrusAdapter{logrusLogger: logrusLogger}
}

func (a *LogrusAdapter) Log(ctx context.Context, level log.Level, keyvals ...interface{}) error {
	var (
		logrusLevel logrus.Level
		fields      logrus.Fields = make(map[string]interface{})
		msg         string
	)

	switch level {
	case log.DEBUG:
		logrusLevel = logrus.DebugLevel
	case log.INFO:
		logrusLevel = logrus.InfoLevel
	case log.WARN:
		logrusLevel = logrus.WarnLevel
	case log.ERROR:
		logrusLevel = logrus.ErrorLevel
	case log.FATAL:
		logrusLevel = logrus.FatalLevel
	default:
		logrusLevel = logrus.DebugLevel
	}

	if logrusLevel > a.logrusLogger.Level {
		return nil
	}

	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		fields[key] = keyvals[i+1]
	}

	if len(fields) > 0 {
		a.logrusLogger.WithFields(fields).Log(logrusLevel, msg)
	} else {
		a.logrusLogger.Log(logrusLevel, msg)
	}
	return nil
}
