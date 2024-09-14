package zaplogx

import (
	"fmt"
	"os"
	"time"

	"github.com/yu1ec/go-pkg/dirx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level      string `json:"level"`
	File       string `json:"file"`
	Production bool   `json:"production"`
}

var (
	stderr = zapcore.Lock(os.Stderr)
	lvl    = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	l      = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), stderr, lvl))
	s      = l.Sugar()
	nop    = zap.NewNop()
)

func NewLogger(lc LogConfig) (*zap.Logger, error) {
	lvl, err := zapcore.ParseLevel(lc.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var out zapcore.WriteSyncer
	if lf := lc.File; len(lf) > 0 {
		err := dirx.CreateNestedDirFromFilepath(lf, os.ModeDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create log dir: %w", err)
		}

		f, _, err := zap.Open(lf)
		if err != nil {
			return nil, fmt.Errorf("open log file failed: %w", err)
		}
		out = zapcore.Lock(f)
	} else {
		out = stderr
	}

	if lc.Production {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
		}
		encoderConfig.TimeKey = "time"
		return zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), out, lvl)), nil
	}

	return zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), out, lvl)), nil
}

// L returns the global logger.
func L() *zap.Logger {
	return l
}

// SetLevel sets the global logger level.
func SetLevel(l zapcore.Level) {
	lvl.SetLevel(l)
}

// S returns the global sugared logger.
func S() *zap.SugaredLogger {
	return s
}

// Nop is a logger that never writes out logs.
func Nop() *zap.Logger {
	return nop
}
