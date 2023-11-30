package vendor

import (
	"encoding/json"
	"fmt"
	"log"
	"nav_sync/config"
	"nav_sync/logger"
	filesystem "nav_sync/mods/afile_system"
	"nav_sync/mods/amanager"
	navapi "nav_sync/mods/anav_api"
	normalapi "nav_sync/mods/anormal_api"
	data_parser "nav_sync/mods/aparser"

	"nav_sync/utils"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.Vendor.Fetch.URL
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	PENDING_LOG_FILE_PATH := utils.VENDOR_DONE_LOG_FILE_PATH
	PENDING_FAILURE := utils.VENDOR_PENDING_FAILURE
	PENDING_SUCCESS := utils.VENDOR_PENDING_SUCCESS

	//Fetch vendor data
	response, err := amanager.Fetch(FETCH_URL, normalapi.GET)
	if err != nil {
		message := "Failed:Fetch:1 " + err.Error()
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	}
	utils.Console(response)

	//get current timestamp
	timestamp := utils.GetCurrentTime()

	//Save to pending file
	err = filesystem.Save(PENDING_FILE_PATH, timestamp, response)
	if err != nil {
		message := "Failed:Fetch:2 " + err.Error()
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, timestamp+".json", message, "")
	} else {
		message := "Fetch: Successfully saved vendor to pending file"
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_SUCCESS, timestamp+".json", message, "")
	}

}

// func Sync() {
// 	//Eg.
// 	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
// 	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
// 	url := config.Config.Vendor.Sync.URL
// 	xmlPayload := `
// 		<Envelope
// 			xmlns="http://schemas.xmlsoap.org/soap/envelope/">
// 			<Body>
// 				<Create
// 					xmlns="urn:microsoft-dynamics-schemas/page/wsvendor">
// 					<WSVendor>
// 						<Name>Suman Intuji </Name>
// 						<Address>From vs code</Address>
// 						<Weighbridge_Supplier_ID>INJ123</Weighbridge_Supplier_ID>
// 					</WSVendor>
// 				</Create>
// 			</Body>
// 		</Envelope>
// 	`

// 	utils.Console(xmlPayload)
// 	result, err := amanager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
// 	if err != nil {
// 		utils.Console(err)
// 	} else {
// 		utils.Console(result)
// 	}
// }

func Sync() {
	//Path
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.VENDOR_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.VENDOR_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(PENDING_FILE_PATH)
	if err != nil {
		message := "Failed:Sync:1 " + err.Error()
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, "", message, "")
	}

	utils.Console(fileNames)

	if fileNames == nil || len(fileNames) < 1 {
		return
	}

	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])

		jsonString := string(jsonData)

		// Unmarshal JSON to struct
		var vendor WSVendor
		if err := json.Unmarshal([]byte(jsonData), &vendor); err != nil {
			message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
		}

		//utils.Console(vendor)

		// Map Go struct to XML
		xmlData, err := data_parser.ParseJsonToXml(vendor)
		if err != nil {
			message := "Failed:Sync:3 Error mapping to XML -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
		}

		//Add XML envelope and body elements
		xmlPayload := fmt.Sprintf(
			`
				<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
					<Body>
						<Create xmlns="urn:microsoft-dynamics-schemas/page/wsvendor">
							%s
						</Create>
					</Body>
				</Envelope>
			`,
			string(xmlData),
		)

		//Return the result
		log.Println(xmlPayload)
		utils.Console("username: ", NTLM_USERNAME)
		utils.Console("username: ", NTLM_PASSWORD)
		utils.Console("URL: ", url)

		//Sync to Nav
		isSuccess := false
		result, err := amanager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
		if err != nil {
			message := "Failed:Sync:4 " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, xmlPayload)
		} else {
			resultStr, ok := result.(string)
			if !ok {
				// The type assertion failed
				message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, xmlPayload)
			}
			match := utils.MatchRegexExpression(resultStr, `<Create_Result[^>]*>`)
			matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

			// Print the result
			if !match && matchFault {
				message := fmt.Sprintf("Failed:Sync:6 XML string does not contain <Create_Result> element: ", result)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
			} else {
				isSuccess = true
			}
		}

		if isSuccess {
			//Move to done file
			err = filesystem.MoveFile(fileNames[i], PENDING_FILE_PATH, DONE_FILE_PATH)
			if err != nil {
				message := "Failed:Sync:5 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, result.(string))
			} else {
				message := "File moved successfully"
				utils.Console(message)
				logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_SUCCESS, fileNames[i], message, "")
			}
		}
	}
}
