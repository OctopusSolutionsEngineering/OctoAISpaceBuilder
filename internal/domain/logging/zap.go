package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ConfigureZapLogger() {
	// Create a custom encoder configuration for plain text logging
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 for timestamps
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"

	// Create a core with a console encoder for plain text output
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // Plain text encoder
		zapcore.Lock(os.Stdout),                  // Output to stdout
		zapcore.InfoLevel,                        // Minimum log level
	)

	// Replace the global logger with the configured logger
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}
