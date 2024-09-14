package zaplogx

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	t.Run("默认配置", func(t *testing.T) {
		lc := LogConfig{
			Level: "info",
		}
		logger, err := NewLogger(lc)
		require.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("生产环境配置", func(t *testing.T) {
		lc := LogConfig{
			Level:      "info",
			Production: true,
		}
		logger, err := NewLogger(lc)
		require.NoError(t, err)
		assert.NotNil(t, logger)

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
		}
		encoderConfig.TimeKey = "time"

		// 测试日志输出
		buf := &bytes.Buffer{}
		testLogger := logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(buf),
				zapcore.InfoLevel,
			)
		}))

		testLogger.Info("test message")

		var logEntry map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &logEntry)
		require.NoError(t, err)

		assert.Equal(t, "test message", logEntry["msg"])
		assert.Contains(t, logEntry["time"].(string), time.Now().Format("2006-01-02 15:04:05.000"))
	})

	t.Run("文件输出", func(t *testing.T) {
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test.log")

		lc := LogConfig{
			Level: "debug",
			File:  logFile,
		}
		logger, err := NewLogger(lc)
		require.NoError(t, err)
		assert.NotNil(t, logger)

		logger.Debug("test debug message")

		content, err := os.ReadFile(logFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), "test debug message")
	})

	t.Run("无效日志级别", func(t *testing.T) {
		lc := LogConfig{
			Level: "invalid",
		}
		_, err := NewLogger(lc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})
}

func TestSetLevel(t *testing.T) {
	originalLevel := lvl.Level()
	defer SetLevel(originalLevel)

	SetLevel(zapcore.DebugLevel)
	assert.Equal(t, zapcore.DebugLevel, lvl.Level())

	SetLevel(zapcore.ErrorLevel)
	assert.Equal(t, zapcore.ErrorLevel, lvl.Level())
}

func TestL(t *testing.T) {
	logger := L()
	assert.NotNil(t, logger)

	buf := &bytes.Buffer{}
	testLogger := logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(buf),
			zapcore.InfoLevel,
		)
	}))

	testLogger.Info("test global logger")
	assert.True(t, strings.Contains(buf.String(), "test global logger"))
}

func TestS(t *testing.T) {
	sugar := S()
	assert.NotNil(t, sugar)

	buf := &bytes.Buffer{}
	testSugar := sugar.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(buf),
			zapcore.InfoLevel,
		)
	}))

	testSugar.Infof("test %s logger", "sugar")
	assert.True(t, strings.Contains(buf.String(), "test sugar logger"))
}

func TestNop(t *testing.T) {
	nopLogger := Nop()
	assert.NotNil(t, nopLogger)

	// 使用 zap.NewNop() 创建一个已知的 NopLogger
	expectedNopLogger := zap.NewNop()

	// 比较我们的 Nop() 返回的 logger 和 zap.NewNop()
	assert.Equal(t, expectedNopLogger, nopLogger)

	// 验证 NopLogger 的核心是 zapcore.NewNopCore()
	assert.Equal(t, zapcore.NewNopCore(), nopLogger.Core())

	// 尝试记录一些信息（这不会产生任何输出，但也不会引发错误）
	nopLogger.Info("this should not be logged")
	nopLogger.Error("this also should not be logged")

	// 无需检查输出，因为 NopLogger 保证不会产生任何输出
}
