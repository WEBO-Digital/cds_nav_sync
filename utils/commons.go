package utils

import (
	"fmt"
	"log"
)

func Console(message interface{}) {
	fmt.Println(message)
}

func Fatal(message interface{}) {
	log.Fatal(message)
}
