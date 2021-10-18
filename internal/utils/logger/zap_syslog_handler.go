package logger

import "go.uber.org/zap/zapcore"

// ZapLoggerHookHandler
// A single log is a structure. This handler intercept each log, you can do any post operations.
// For example: send log to elasticsearch
// Entry details:
// Level: log level
// Time: current timestamp
// LoggerName: log name
// Message: log content
// Caller: each file callers' path
// Stack: code stack
func ZapLoggerHookHandler(entry zapcore.Entry) error {
	go func(e zapcore.Entry) {
		// fmt.Println("Benlog hook, you can continue to handle system log here...")
	}(entry)
	return nil
}
