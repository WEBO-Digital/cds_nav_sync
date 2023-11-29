package logger

import (
	"fmt"
	filesystem "nav_sync/mods/afile_system"
	"nav_sync/utils"
)

// Creating different types of response
type ResponseType string

const (
	SUCCESS ResponseType = "SUCCESS"
	FAILURE ResponseType = "FAILURE"
)

func LogInvoiceFetch(responseType ResponseType, filePath string, fileName string, savedfileName string, message string, data interface{}) {
	// Type assertion to convert interface to string
	str, ok := data.(string)
	if ok {
		// Successfully converted to string
		fmt.Println("String:", str)
	} else {
		// Conversion failed
		fmt.Println("Not a string")
	}

	//get current timestamp
	timestamp := utils.GetCurrentTime()

	//Format data: Please do not change its format
	appendStr := fmt.Sprintf(
		`

**********************************START*************************************
[%s]
Type: %s
File Name: %s
Message: %s
**********************************END***************************************

		`, timestamp, responseType, savedfileName, message,
	)

	//Save to particular path
	err := filesystem.Append(filePath, fileName, appendStr)
	if err != nil {
		utils.Console(err.Error())
	}
}
