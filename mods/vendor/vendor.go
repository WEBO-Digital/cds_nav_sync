package vendor

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"nav_sync/config"
	"nav_sync/logger"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/mods/ahelpers/manager"
	navapi "nav_sync/mods/ahelpers/nav_api"
	normalapi "nav_sync/mods/ahelpers/normal_api"
	data_parser "nav_sync/mods/ahelpers/parser"
	"nav_sync/utils"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.Vendor.Fetch.URL
	TOKEN_KEY := config.Config.Vendor.Fetch.APIKey
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	PENDING_LOG_FILE_PATH := utils.VENDOR_DONE_LOG_FILE_PATH
	PENDING_FAILURE := utils.VENDOR_PENDING_FAILURE
	PENDING_SUCCESS := utils.VENDOR_PENDING_SUCCESS

	//Fetch vendor data
	response, err := manager.Fetch(FETCH_URL, normalapi.GET, TOKEN_KEY, nil)
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

	var responseModel []BackToCDSVendorResponse
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
		result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
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

		isSuccessfullySavedToFile := false
		if isSuccess {
			//Move to done file
			err = filesystem.MoveFile(fileNames[i], PENDING_FILE_PATH, DONE_FILE_PATH)
			if err != nil {
				message := "Failed:Sync:5 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, result.(string))
			} else {
				isSuccessfullySavedToFile = true
				message := "File moved successfully"
				utils.Console(message)
				logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_SUCCESS, fileNames[i], message, "")
			}
		}

		//Add successed to an array
		if isSuccessfullySavedToFile {
			// Convert the string to a byte slice
			xmlData := []byte(result.(string))

			// Map Go struct to XML
			var parseModel CreateResultVendor
			err := xml.Unmarshal(xmlData, &parseModel)
			if err != nil {
				message := "Failed:Sync:6 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, result.(string))
			}

			// responseModel[i].VendorNo = parseModel.Body.CreateResult.WSVendor.No
			// responseModel[i].WeighbridgeSupplierID = parseModel.Body.CreateResult.WSVendor.WeighbridgeSupplierID

			responseModel = append(responseModel, BackToCDSVendorResponse{
				VendorNo:              parseModel.Body.CreateResult.WSVendor.No,
				WeighbridgeSupplierID: parseModel.Body.CreateResult.WSVendor.WeighbridgeSupplierID,
			})
		}
	}

	//Bulk save
	//After syncing all files then send response back to CDS
	utils.Console("responseModel------> ", len(responseModel))
	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func sendToCDS(responseModel []BackToCDSVendorResponse) {
	//Path
	RESPONSE_URL := config.Config.Vendor.Save.URL
	TOKEN_KEY := config.Config.Vendor.Fetch.APIKey

	//Save Response vendor data to CDS
	response, err := manager.Fetch(RESPONSE_URL, normalapi.POST, TOKEN_KEY, responseModel)
	if err != nil {
		message := "Failed:sendToCDS:Fetch:1 " + err.Error()
		utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	}
	if response != nil {
		utils.Console(response)
	}

	// utils.Console("Successfully send to CDS system from nav ---> vendor: ", RESPONSE_URL)
	// utils.Console(responseModel)
}

func Sync2() {
	//Path
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.VENDOR_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.VENDOR_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS

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

	//Get Hash Database
	hashModels := GetHash()

	//Syncing
	var responseModel []BackToCDSVendorResponse
	for i := 0; i < len(fileNames); i++ {
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])
		jsonString := string(jsonData)

		// Unmarshal JSON to struct
		var vendors []WSVendor
		if err := json.Unmarshal([]byte(jsonData), &vendors); err != nil {
			message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
		}

		var isSuccessArr []bool
		for j := 0; j < len(vendors); j++ {
			hashMaps, isHashed := CompareWithHash(vendors[j], hashModels)
			if !isHashed {
				isSuccess, err, result := InsertToNav(vendors[j], fileNames[i], jsonString)
				if err != nil {
					message := err.Error() //"Failed:Sync:3 Error mapping to XML -> " + err.Error()
					utils.Console(message)
					logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
				}
				isSuccessArr = append(isSuccessArr, isSuccess)
				if isSuccess {
					// Convert the string to a byte slice
					xmlData := []byte(result.(string))

					// Map Go struct to XML
					var parseModel CreateResultVendor
					err = xml.Unmarshal(xmlData, &parseModel)
					if err != nil {
						message := "Failed:Sync:6 " + err.Error()
						utils.Console(message)
						logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, result.(string))
					}

					//append success hash map and save hash map
					vendorStr, _ := data_parser.ParseModelToString(vendors[j])

					// Update the Hash field for a specific key
					err = UpdateHashInModel(hashMaps, vendors[j].WeighbridgeSupplierID, utils.ComputeMD5(vendorStr), parseModel.Body.CreateResult.WSVendor.No)
					if err != nil {
						utils.Console("UpdateHashInModel::Error:", err)
					}
					utils.Console("hashMaps", hashMaps)

					//Add successed to an array
					responseModel = append(responseModel, BackToCDSVendorResponse{
						VendorNo:              parseModel.Body.CreateResult.WSVendor.No,
						WeighbridgeSupplierID: parseModel.Body.CreateResult.WSVendor.WeighbridgeSupplierID,
					})

				}

			}
		}
		//Check if can be saved to file system
		// canSaveToFileSystem := true
		// for j := 0; j < len(isSuccessArr); j++ {
		// 	if !isSuccessArr[j] {
		// 		canSaveToFileSystem = false
		// 		break
		// 	}
		// }

		//Move to done file
		//if canSaveToFileSystem {
		err = filesystem.MoveFile(fileNames[i], PENDING_FILE_PATH, DONE_FILE_PATH)
		if err != nil {
			message := "Failed:Sync:5 " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
		} else {
			//isSuccessfullySavedToFile = true
			message := "File moved successfully"
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_SUCCESS, fileNames[i], message, "")
		}
		//}
	}

	//Save to Hash Folder
	val, err := SaveHashLogs(hashModels)
	if err != nil {
		utils.Console("SaveHashLogs::err", err)
	}
	utils.Console("SaveHashLogs::val", val)

	//Bulk save
	//After syncing all files then send response back to CDS
	utils.Console("responseModel------> ", len(responseModel))
	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}
