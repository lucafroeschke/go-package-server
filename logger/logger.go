package logger

import (
	"fmt"
	"time"
)

const (
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
)

func WriteLog(logType string, message string) {
	timestamp := time.Now().Format(time.RFC3339)

	fmt.Println(
		fmt.Sprintf("[%s] [%s] %s", timestamp, logType, message),
	)
}
