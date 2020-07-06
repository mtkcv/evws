package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

//LoggerInit init logger
func LoggerInit(client bool) {

	logPath := "evws.log"
	if client {
		logPath = "evws_client.log"
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		logPath,
	}
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		CallerKey:   "caller",
		MessageKey:  "content",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	build, err := cfg.Build()

	if err != nil {
		fmt.Printf("zap log build error, %s\n", err)
		return
	}
	logger = build.Sugar()
}
