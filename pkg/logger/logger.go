package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type (
	LogHandler interface {
		Info(args ...any)
		Infof(template string, args ...any)
		Debug(args ...any)
		Debugf(template string, args ...any)
		Error(args ...any)
		Errorf(template string, args ...any)
		Fatalf(template string, args ...any)
		PrettifyJSON(v any)
	}

	Logger struct {
		*zap.SugaredLogger
	}
)

func New() *Logger {
	var logger *zap.Logger
	logger, _ = zap.NewProduction()
	sugarLogger := logger.Sugar()
	return &Logger{sugarLogger}
}

func (l *Logger) PrettifyJSON(v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("error marshaling to JSON: %v\n", err)
	}
	fmt.Println(string(jsonData))
}
