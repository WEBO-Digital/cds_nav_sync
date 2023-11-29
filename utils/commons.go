package utils

import (
	"fmt"
	"log"
	"regexp"
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

func MatchRegexExpression(value string, pattern string) bool {
	// Define the regular expression pattern
	//pattern := `<Create_Result[^>]*>`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the pattern matches the XML string
	match := re.MatchString(value)
	return match
}
