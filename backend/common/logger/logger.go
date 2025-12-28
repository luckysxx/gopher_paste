package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 创建一个新的 Logger 实例
// serviceName 用于标识日志来源的服务名称
func NewLogger(serviceName string) *zap.Logger {
	config := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",

		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法打开日志文件:", err)
		return zap.NewExample()
	}

	// 开发环境使用Debug级别，生产环境使用Info级别
	level := zapcore.InfoLevel
	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "development" {
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(file),
		),
		level,
	)

	// AddCaller添加调用者信息，AddStacktrace在Error级别添加堆栈
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(0),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	logger = logger.With(zap.String("service", serviceName))

	return logger
}
