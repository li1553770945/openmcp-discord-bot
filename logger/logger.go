package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(logDir string) {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}

	// 创建两个独立encoder（关键改动点）
	consoleEncoder := getColorfulEncoder() // 控制台：带颜色
	fileEncoder := getPlainTextEncoder()   // 文件：无颜色

	// 两个不同的输出目标
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileWriter := getLogWriter(filepath.Join(logDir, "app.log"))

	// 分别创建Core
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, zapcore.InfoLevel)

	// 合并Core
	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core, zap.AddCaller())

	// 替换全局实例
	zap.ReplaceGlobals(logger)
}

// 控制台编码器（带颜色）
func getColorfulEncoder() zapcore.Encoder {
	conf := zap.NewDevelopmentEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder // 关键颜色配置
	return zapcore.NewConsoleEncoder(conf)
}

// 文件编码器（无颜色纯文本）
func getPlainTextEncoder() zapcore.Encoder {
	conf := zap.NewDevelopmentEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = zapcore.CapitalLevelEncoder // 关键去颜色配置
	return zapcore.NewConsoleEncoder(conf)
}

func getLogWriter(logPath string) zapcore.WriteSyncer {
	rotator, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(rotator)
}
