package logger

import (
	"microservice/pkg/env"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Colores ANSI
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m" //Error
	ColorYellow = "\033[33m" //Warn
	ColorWhite  = "\033[37m" //Info
	ColorGreen  = "\033[32m" //Success
	ColorPurple = "\033[35m" //General/Debug
)

var (
	log  *zap.Logger
	once sync.Once // Singleton
)

func Init(environment string) {
	once.Do(func() {

		env := env.Parse(environment)

		var cfg zap.Config

		if env.IsProduction() {
			cfg = zap.NewProductionConfig()
			cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		} else {
			core := zapcore.NewCore(
				createColoredEncoder(),
				zapcore.AddSync(os.Stdout),
				zapcore.DebugLevel,
			)
			log = zap.New(core, zap.AddCaller())
			return
		}

		l, err := cfg.Build()
		if err != nil {
			fallback := zap.NewExample()
			fallback.Error("No se pudo inicializar zap correctamente, usando fallback", zap.Error(err))
			log = fallback
			return
		}

		log = l
	})
}

func L() *zap.Logger {
	if log == nil {
		panic("logger no inicializado — llama a logger.Init(env) antes de usar logger.L()")
	}
	return log
}

func Sugar() *zap.SugaredLogger {
	return L().Sugar()
}

func Success(msg string, fields ...zap.Field) {
	L().Info(ColorGreen+"✓ "+msg+ColorReset, fields...)
}

func General(msg string, fields ...zap.Field) {
	L().Debug(ColorWhite+"✓ "+msg+ColorReset, fields...)
}

func Error(msg string, fields ...zap.Field) {
	L().Error(ColorRed+"✓ "+msg+ColorReset, fields...)
}

func createColoredEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()

	config.EncodeLevel = coloredLevelEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(config)

}

func coloredLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(ColorPurple + "[DEBUG]" + ColorReset)
	case zapcore.InfoLevel:
		enc.AppendString(ColorWhite + "[INFO]" + ColorReset)
	case zapcore.WarnLevel:
		enc.AppendString(ColorYellow + "[WARN]" + ColorReset)
	case zapcore.ErrorLevel:
		enc.AppendString(ColorRed + "[ERROR]" + ColorReset)
	case zapcore.FatalLevel:
		enc.AppendString(ColorRed + "[FATAL]" + ColorReset)
	default:
		enc.AppendString(level.CapitalString())
	}
}
