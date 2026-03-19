// Stub package — the real logger lives in the svw_mono monorepo (gsail-go/logger).
// This stub exists only so stage_primer compiles standalone for code review.
package logger

import "go.uber.org/zap"

func InitGlobalLogger(logFile string, customEncoder interface{}, compress bool) {}

func Get() *zap.SugaredLogger {
	l, _ := zap.NewDevelopment()
	return l.Sugar()
}
