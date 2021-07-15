package logger

import "go.uber.org/zap"

// TODO create general interface with generic fields

func Debug(msg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infof(format, args...)
}

func Warn(msg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Warnf(format, args...)
}

func Error(msg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Errorf(format, args...)
}
