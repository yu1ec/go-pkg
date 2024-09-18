package zaplogx

import (
	"fmt"
	"os"
	"time"

	"github.com/yu1ec/go-pkg/dirx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RotateConfig 日志轮转配置
type RotateConfig struct {
	MaxSize    int  `json:"maxsize" default:"100"`    // 最大文件大小，单位：MB
	MaxAge     int  `json:"maxage" default:"7"`       // 最大保存天数
	MaxBackups int  `json:"maxbackups" default:"10"`  // 最大备份数
	Compress   bool `json:"compress" default:"false"` // 是否压缩
	LocalTime  bool `json:"localtime" default:"true"` // 是否使用本地时间
}

type LogConfig struct {
	Level      string        `json:"level"`
	File       string        `json:"file"`
	Production bool          `json:"production"`
	Rotate     *RotateConfig `json:"rotate,omitempty"`
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

		if lc.Rotate != nil {
			rotateCfg := lc.Rotate
			rotateLog := &lumberjack.Logger{
				Filename:   lc.File,
				MaxSize:    rotateCfg.MaxSize,
				MaxAge:     rotateCfg.MaxAge,
				MaxBackups: rotateCfg.MaxBackups,
				Compress:   rotateCfg.Compress,
			}
			out = zapcore.AddSync(rotateLog)
		} else {
			f, _, err := zap.Open(lf)
			if err != nil {
				return nil, fmt.Errorf("open log file failed: %w", err)
			}
			out = zapcore.Lock(f)
		}
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
