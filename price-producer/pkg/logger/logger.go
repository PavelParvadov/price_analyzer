package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger() *zap.Logger {
	cfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		ConsoleSeparator: " | ",
	}
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger

}
