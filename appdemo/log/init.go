package sugar

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"runtime"
)

var Sugar *zap.SugaredLogger
var output io.WriteCloser

func init() {
	var logger *zap.Logger
	var level zapcore.LevelEnabler

	if runtime.GOOS == "windows" {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}
	writer := getLogWriter()

	core := zapcore.NewCore(getEncoder(), writer, level)
	logger = zap.New(core, zap.AddCaller())
	Sugar = logger.Sugar()

	// 测试日志分割
	//for i := 0; i < 10000000; i++ {
	//	Sugar.Warnln(longStr())
	//}
}

var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func longStr() string {
	s1 := s
	for i := 0; i < 10; i++ {
		s1 += s
	}
	return s1
}

func Stop() {
	Sugar.Sync()
	output.Close()
}

func getEncoder() zapcore.Encoder {
	var encoderConf zapcore.EncoderConfig
	if runtime.GOOS == "windows" {
		encoderConf = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConf = zap.NewProductionEncoderConfig()
	}
	return zapcore.NewConsoleEncoder(encoderConf)
}

func getLogWriter() zapcore.WriteSyncer {
	var output io.WriteCloser
	var file = "project.log"
	if runtime.GOOS == "windows" {
		output = os.Stdout
	} else {
		// 日志分割
		output = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    100, // 单文件MB
			MaxBackups: 5,   // 保留日志文件的数量
			MaxAge:     30,  // 保留30天日志
			Compress:   false,
		}
	}
	return zapcore.AddSync(output)
}

// --------------------------------

func Debug(args ...interface{}) {
	Sugar.Debugln(args...)
}

func Info(args ...interface{}) {
	Sugar.Infoln(args...)
}
func Infof(t string, args ...interface{}) {
	Sugar.Infof(t, args...)
}

func Warn(args ...interface{}) {
	Sugar.Warnln(args...)
}

func Error(args ...interface{}) {
	Sugar.Errorln(args...)
}

func Errorf(t string, args ...interface{}) {
	Sugar.Errorf(t, args...)
}

func Panic(args ...interface{}) {
	Sugar.Panicln(args...)
}
