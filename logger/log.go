package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

/* -------------------- Exported Functions -------------------- */

func init() {
	Log("Logging initialized")
}

var once sync.Once

func Log(msg string) {
	if LogFileMissing() {
		return
	}

	once.Do(func() {
		f, err := os.OpenFile(LogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0600)

		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		log.SetOutput(f)
	})

	log.Println(msg)
}

func LogJson(message string, data any) {
	dataStr, _ := json.MarshalIndent(data, "  ", "  ")

	Log(fmt.Sprintf("%s: %s", message, dataStr))
}

func LogFileMissing() bool {
	return LogFilePath() == ""
}

func LogFilePath() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, ".config", "wtf", "log.txt")
}
