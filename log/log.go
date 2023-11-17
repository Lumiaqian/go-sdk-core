package log

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

type LogHelper struct {
	logger Logger
}

func NewLogHelper(logger Logger) *LogHelper {
	return &LogHelper{
		logger: logger,
	}
}

func (l *LogHelper) Debug(keyvals ...interface{}) {
	l.logger.Log(DEBUG, keyvals...)
}

func (l *LogHelper) Info(keyvals ...interface{}) {
	l.logger.Log(INFO, keyvals...)
}

func (l *LogHelper) Warn(keyvals ...interface{}) {
	l.logger.Log(WARN, keyvals...)
}

func (l *LogHelper) Error(keyvals ...interface{}) {
	l.logger.Log(ERROR, keyvals...)
}
