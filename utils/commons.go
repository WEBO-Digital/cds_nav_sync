package utils

import (
	"fmt"
	"log"
	"time"
)

func Console(message ...any) {
	fmt.Println(message)
}

func Fatal(message ...any) {
	log.Fatal(message)
}

func GetCurrentTime() string {
	timestamp := time.Now().Format("2006-01-02T15-04-05.999")
	return timestamp
}
