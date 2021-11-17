package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger interface {
	Info(message string, args ...zapcore.Field)
	Error(message string, err error, args ...zapcore.Field)
	Fatal(message string, err error, args ...zapcore.Field)
	With(arg ...zapcore.Field) AppLogger
}

type appLogger struct {
	l *zap.Logger
}

var aLogger appLogger

func InitLogger(appName string) (*appLogger, error) {
	cnf := zap.NewProductionConfig()
	cnf.DisableStacktrace = true
	cnf.DisableCaller = true
	cnf.EncoderConfig.TimeKey = zapcore.OmitKey

	z, err := cnf.Build()
	if err != nil {
		return nil, err
	}
	aLogger = appLogger{
		l: z.With(zap.String("app", appName)),
	}
	return &aLogger, nil
}

func Info(message string, args ...zapcore.Field) {
	aLogger.l.Info(message, args...)
}

func Error(message string, err error, args ...zapcore.Field) {
	if len(args) == 0 {
		aLogger.l.Error(message, zap.Error(err))
		return
	}
	aLogger.l.Error(message, prepareParams(err, args)...)
}

func Fatal(message string, err error, args ...zapcore.Field) {
	if len(args) == 0 {
		aLogger.l.Fatal(message, zap.Error(err))
		return
	}
	aLogger.l.Fatal(message, prepareParams(err, args)...)
}

func prepareParams(err error, args []zapcore.Field) []zapcore.Field {
	params := make([]zapcore.Field, 0, len(args)+1)
	params = append(params, zap.Error(err))
	params = append(params, args...)
	return params
}
