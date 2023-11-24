package utils

import (
	"fmt"
	"log"
)

func Console(message ...any) {
	fmt.Println(message)
}

func Fatal(message ...any) {
	log.Fatal(message)
}
