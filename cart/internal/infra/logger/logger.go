package logger

import (
	"fmt"
	"time"
)

func log(level string, message string) {
	fmt.Printf("[%s] %s - %s\n",
		level, time.Now().Format("2006/01/02 15:04:05"), message)
}

// Выводит информационное сообщение в лог
func Info(message string) {
	log("INFO", message)
}
