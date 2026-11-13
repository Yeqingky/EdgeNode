package kvstore

import (
	"fmt"

	"git.ovh.im/yeqing/EdgeNode/internal/remotelogs"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (this *Logger) Infof(format string, args ...any) {
	// stub
}

func (this *Logger) Errorf(format string, args ...any) {
	remotelogs.Error("KV", fmt.Sprintf(format, args...))
}

func (this *Logger) Fatalf(format string, args ...any) {
	remotelogs.Error("KV", fmt.Sprintf(format, args...))
}
