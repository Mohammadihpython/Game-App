package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var Logger *zap.Logger

var once = sync.Once{}

func init() {
	// ما میدونیم تابع init یکبار اچرا میشه ولی از once تو سناریو های دیگه  برای اینکه یبار صدا زده بشه استفاده میکنیم
	once.Do(func() {
		Logger, _ = zap.NewDevelopment()

	})
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	//logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//fileWriter := zapcore.AddSync(logFile)
	// Add lumberjack log rotate ability
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/log.json",
		MaxSize:    2, // megabytes
		MaxBackups: 3,
		LocalTime:  false,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	stdOutWriter := zapcore.AddSync(os.Stdout)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel),
		zapcore.NewCore(fileEncoder, stdOutWriter, zap.InfoLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

}
