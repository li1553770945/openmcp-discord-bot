package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

func InitLogger(logDir string) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}
	// 配置文本格式编码器
	encoder := getTextEncoder()
	// 设置两个输出目标（控制台 + 滚动文件）
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileWriter := getLogWriter(filepath.Join(logDir, "app.log"))
	// 创建两个Core（核心）
	// 控制台：调试级别+全部显示，文件：信息级别+无StackTrace
	consoleCore := zapcore.NewCore(encoder, consoleWriter, zapcore.DebugLevel)
	fileCore := zapcore.NewCore(encoder, fileWriter, zapcore.InfoLevel)
	// 合并Core，创建Logger
	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	// 替换Zap全局Logger和SugarLogger
	zap.ReplaceGlobals(logger)
}
func getTextEncoder() zapcore.Encoder {
	// 纯文本格式（带颜色）
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色日志级别
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func getLogWriter(logPath string) zapcore.WriteSyncer {
	// 滚动日志配置（按天分割，保留7天）
	rotator, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),      // 生成软链
		rotatelogs.WithMaxAge(7*24*time.Hour), // 保留7天
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(rotator)
}
