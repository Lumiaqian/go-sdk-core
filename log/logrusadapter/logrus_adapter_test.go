package logrusadapter

import (
	"bytes"
	"testing"

	"github.com/Lumiaqian/go-sdk-core/log"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogrusAdapter_Log(t *testing.T) {
	// Create a logrus instance with a buffer to read and verify the output
	var buf bytes.Buffer
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(&buf)
	logrusLogger.SetLevel(logrus.DebugLevel)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})

	adapter := NewLogrusAdapter(logrusLogger)

	// Test different log levels
	testCases := []struct {
		level    log.Level
		keyvals  []interface{}
		expected string
		message  string
	}{
		{log.DEBUG, []interface{}{"message", "debug message"}, "debug message", "Debug level test"},
		{log.INFO, []interface{}{"message", "info message"}, "info message", "Info level test"},
		{log.WARN, []interface{}{"message", "warn message"}, "warn message", "Warn level test"},
		{log.ERROR, []interface{}{"message", "error message"}, "error message", "Error level test"},
		{log.FATAL, []interface{}{"message", "fatal message"}, "fatal message", "Fatal level test"},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			buf.Reset()
			err := adapter.Log(tc.level, tc.keyvals...)
			assert.NoError(t, err)
			if tc.expected != "" {
				assert.Contains(t, buf.String(), tc.expected)
			}
		})
	}

	// Test invalid key (non-string key)
	t.Run("Invalid key", func(t *testing.T) {
		buf.Reset()
		adapter.Log(log.INFO, 123, "invalid key")
		assert.NotContains(t, buf.String(), "invalid key")
	})

	// Test odd number of keyvals
	t.Run("Odd keyvals", func(t *testing.T) {
		buf.Reset()
		adapter.Log(log.INFO, "odd", "number", "of", "keyvals", "key")
		assert.Contains(t, buf.String(), "\"key\":\"\"") // Last key has empty value
		assert.Contains(t, buf.String(), "\"odd\":\"number\"")
		assert.Contains(t, buf.String(), "\"of\":\"keyvals\"")
	})
}
