package logger

import (
	"fmt"
	"os"
	"time"
)

func LogError(message string) {
	f, err := os.OpenFile("statusapp_errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		// If we can't log, silently ignore to avoid impacting the application
		return
	}
	defer f.Close()
	ts := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s - %s\n", ts, message)
	_, _ = f.WriteString(line)
}
