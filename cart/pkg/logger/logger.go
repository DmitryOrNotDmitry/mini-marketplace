package logger

import (
	"fmt"
	"time"
)

func log(level string, message string) {
	fmt.Printf("[%s] %s - %s\n",
		level, time.Now().Format("2006/01/02 15:04:05"), message)
}

// Info выводит информационное сообщение в лог
func Info(message string) {
	log("INFO", message)
}

// Error выводит сообщение об ошибке в лог
func Error(message string) {
	log("ERROR", message)
}

// Warning выводит сообщение об предупреждении в лог
func Warning(message string) {
	log("WARNING", message)
}
