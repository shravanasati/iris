package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	logger  *log.Logger
	logFile *os.File
)

// InitLogger initializes the logger to write to a daily file in ~/.iris/logs
func InitLogger() {
	logDir := filepath.Join(GetIrisDir(), "logs")

	// Create current day log file
	currentTime := time.Now()
	fileName := fmt.Sprintf("%s.log", currentTime.Format("2006-01-02"))
	logPath := filepath.Join(logDir, fileName)

	var err error
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		return
	}

	logger = log.New(logFile, "", log.LstdFlags)
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// logMessage is internal helper to format and write logs
func logMessage(level string, component string, message string) {
	if logger == nil {
		return
	}

	// format: [COMPONENT] level: message
	msg := fmt.Sprintf("[%s] %s: %s", strings.ToLower(component), strings.ToLower(level), strings.ToLower(message))
	logger.Println(msg)
}

func LogInfof(component string, format string, v ...any) {
	logMessage("INFO", component, fmt.Sprintf(format, v...))
}

func LogErrorf(component string, format string, v ...any) {
	logMessage("ERROR", component, fmt.Sprintf(format, v...))
}

func LogWarnf(component string, format string, v ...any) {
	logMessage("WARN", component, fmt.Sprintf(format, v...))
}

// CleanupLogs removes log files older than 31 days
func CleanupLogs() {
	logDir := filepath.Join(GetIrisDir(), "logs")
	files, err := os.ReadDir(logDir)
	if err != nil {
		return
	}

	now := time.Now()
	retentionPeriod := 31 * 24 * time.Hour

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > retentionPeriod {
			os.Remove(filepath.Join(logDir, file.Name()))
			LogInfof("cleanup", "deleted old log file: %s", file.Name())
		}
	}
}
