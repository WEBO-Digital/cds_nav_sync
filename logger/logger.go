package logger

import (
	"fmt"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/utils"
)

// Creating different types of response
type ResponseType string

const (
	SUCCESS ResponseType = "SUCCESS"
	FAILURE ResponseType = "FAILURE"
)

func LogNavState(responseType ResponseType, filePath string, fileName string, savedfileName string, message string, data string) { //data interface{}) {
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
Data: %s
**********************************END***************************************

		`, timestamp, responseType, savedfileName, message, data,
	)

	//Save to particular path
	err := filesystem.Append(filePath, fileName, appendStr)
	if err != nil {
		utils.Console(err.Error())
	}
}

func AddToLog(filePath string, fileName string, rType ResponseType, message string, data string) {
	timestamp := utils.GetCurrentTime()

	appendStr := fmt.Sprintf(
		`**********************************START*************************************
[%s]
Type: %s
Message: %s
Data: %s
**********************************END***************************************

`, timestamp, rType, message, data,
	)

	err := filesystem.Append(filePath, fileName, appendStr)

	if err != nil {
		utils.Console(err.Error())
	}
}
