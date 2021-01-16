package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func init() {
	//initialise our global logger

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevel(),
	),
		zap.AddStacktrace(zap.FatalLevel), //disable stacktrace for level lower than fatal
		zap.WithCaller(false),
	)
	zap.ReplaceGlobals(logger)

	MQTT.ERROR, _ = zap.NewStdLogAt(zap.L(), zap.ErrorLevel)
	MQTT.CRITICAL, _ = zap.NewStdLogAt(zap.L(), zap.ErrorLevel)
	MQTT.WARN, _ = zap.NewStdLogAt(zap.L(), zap.WarnLevel)
	MQTT.DEBUG, _ = zap.NewStdLogAt(zap.L(), zap.DebugLevel)
}
