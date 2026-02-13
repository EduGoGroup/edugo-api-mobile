package adapter

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoggerAdapter(t *testing.T) {
	// Crear logger de logrus
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)

	// Crear adapter
	adapter := NewLoggerAdapter(logrusLogger)

	// Verificar que no es nil
	require.NotNil(t, adapter)
}

func TestLoggerAdapter_BasicLogging(t *testing.T) {
	// Crear logger de logrus
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)

	// Crear adapter
	adapter := NewLoggerAdapter(logrusLogger)

	// Probar todos los niveles sin fields
	t.Run("Debug without fields", func(t *testing.T) {
		adapter.Debug("test debug message")
	})

	t.Run("Info without fields", func(t *testing.T) {
		adapter.Info("test info message")
	})

	t.Run("Warn without fields", func(t *testing.T) {
		adapter.Warn("test warn message")
	})

	t.Run("Error without fields", func(t *testing.T) {
		adapter.Error("test error message")
	})
}

func TestLoggerAdapter_LoggingWithFields(t *testing.T) {
	// Crear logger de logrus
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)

	// Crear adapter
	adapter := NewLoggerAdapter(logrusLogger)

	// Probar logging con fields
	t.Run("Info with fields", func(t *testing.T) {
		adapter.Info("test message", "key1", "value1", "key2", 123)
	})

	t.Run("Error with fields", func(t *testing.T) {
		adapter.Error("error message", "error", "something went wrong", "code", 500)
	})
}

func TestLoggerAdapter_With(t *testing.T) {
	// Crear logger de logrus
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)

	// Crear adapter
	adapter := NewLoggerAdapter(logrusLogger)

	// Crear logger con contexto
	contextLogger := adapter.With("user_id", "123", "request_id", "abc")
	require.NotNil(t, contextLogger)

	// Usar el logger con contexto
	contextLogger.Info("operation completed")
	contextLogger.Error("operation failed", "error", "timeout")
}

func TestLoggerAdapter_WithChaining(t *testing.T) {
	// Crear logger de logrus
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)

	// Crear adapter
	adapter := NewLoggerAdapter(logrusLogger)

	// Encadenar With
	logger1 := adapter.With("level1", "value1")
	logger2 := logger1.With("level2", "value2")
	logger3 := logger2.With("level3", "value3")

	// Verificar que funciona
	logger3.Info("test chained context")
}

func TestLoggerAdapter_Sync(t *testing.T) {
	// Crear logger de logrus
	logrusLogger := logrus.New()

	// Crear adapter
	adapter := NewLoggerAdapter(logrusLogger)

	// Sync no debería retornar error (logrus no requiere sync)
	err := adapter.Sync()
	assert.NoError(t, err)
}

func TestConvertToLogrusFields(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected logrus.Fields
	}{
		{
			name:     "empty fields",
			input:    []interface{}{},
			expected: logrus.Fields{},
		},
		{
			name:     "single pair",
			input:    []interface{}{"key", "value"},
			expected: logrus.Fields{"key": "value"},
		},
		{
			name:     "multiple pairs",
			input:    []interface{}{"key1", "value1", "key2", 123, "key3", true},
			expected: logrus.Fields{"key1": "value1", "key2": 123, "key3": true},
		},
		{
			name:     "odd number of fields (last ignored)",
			input:    []interface{}{"key1", "value1", "orphan"},
			expected: logrus.Fields{"key1": "value1"},
		},
		{
			name:     "non-string key",
			input:    []interface{}{123, "value"},
			expected: logrus.Fields{"unknown": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToLogrusFields(tt.input...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLoggerAdapter_Debug_WithFields verifica cobertura de Debug con fields
func TestLoggerAdapter_Debug_WithFields(t *testing.T) {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)
	adapter := NewLoggerAdapter(logrusLogger)

	adapter.Debug("debug with fields", "key", "value")
}

// TestLoggerAdapter_Warn_WithFields verifica cobertura de Warn con fields
func TestLoggerAdapter_Warn_WithFields(t *testing.T) {
	logrusLogger := logrus.New()
	adapter := NewLoggerAdapter(logrusLogger)

	adapter.Warn("warn with fields", "key", "value")
}

// TestLoggerEntryAdapter_Debug verifica Debug en entry adapter
func TestLoggerEntryAdapter_Debug(t *testing.T) {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel)
	adapter := NewLoggerAdapter(logrusLogger)

	entryLogger := adapter.With("context", "test")
	entryLogger.Debug("entry debug")
}

// TestLoggerEntryAdapter_Warn verifica Warn en entry adapter
func TestLoggerEntryAdapter_Warn(t *testing.T) {
	logrusLogger := logrus.New()
	adapter := NewLoggerAdapter(logrusLogger)

	entryLogger := adapter.With("context", "test")
	entryLogger.Warn("entry warn")
}

// TestLoggerEntryAdapter_Sync verifica Sync en entry adapter
func TestLoggerEntryAdapter_Sync(t *testing.T) {
	logrusLogger := logrus.New()
	adapter := NewLoggerAdapter(logrusLogger)

	entryLogger := adapter.With("context", "test")
	err := entryLogger.Sync()

	assert.NoError(t, err)
}
