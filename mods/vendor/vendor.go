package vendor

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"nav_sync/config"
	"nav_sync/logger"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/mods/ahelpers/manager"
	normalapi "nav_sync/mods/ahelpers/normal_api"
	data_parser "nav_sync/mods/ahelpers/parser"
	hashrecs "nav_sync/mods/hashrecs"
	"nav_sync/utils"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.Vendor.Fetch.URL
	TOKEN_KEY := config.Config.Vendor.Fetch.APIKey
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	VENDOR_LOG_PATH := utils.VENDOR_LOG_PATH

	utils.Console("Start fetching vendors")

	//get current timestamp
	timestamp := utils.GetCurrentTime()
	logFileName := "fetch-" + timestamp + ".log"

	//Fetch vendor data
	response, err := manager.Fetch(FETCH_URL, normalapi.GET, TOKEN_KEY, nil)

	if err != nil {
		message := "Failed[1]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.FAILURE, message, FETCH_URL)
		return
	}

	//Save to pending file
	err = filesystem.Save(PENDING_FILE_PATH, timestamp, response)

	if err != nil {
		message := "Failed[2]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.FAILURE, message, "")
	} else {
		message := "fetched vendors and saved to a file"
		utils.Console(message)
		logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.SUCCESS, message, PENDING_FILE_PATH+timestamp+".json")
	}

}

func Sync3() {
	//Path
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.VENDOR_DONE_FILE_PATH
	VENDOR_LOG_PATH := utils.VENDOR_LOG_PATH
	HASH_FILE_PATH := utils.VENDOR_HASH_FILE_PATH
	HASH_DB := utils.VENDOR_HASH_DB
	PREFIX := config.Config.Vendor.Prefix

	utils.Console("Start sending vendors to NAV")

	timestamp := utils.GetCurrentTime()
	logFileGeneral := "sync-pre-" + timestamp + ".json"

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(PENDING_FILE_PATH)

	if err != nil {
		message := "Failed[1]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(VENDOR_LOG_PATH, logFileGeneral, logger.FAILURE, message, "")
		return
	}

	if fileNames == nil || len(fileNames) < 1 {
		message := "No pending files found"
		utils.Console(message)
		logger.AddToLog(VENDOR_LOG_PATH, logFileGeneral, logger.SUCCESS, message, "")
		return
	}

	//Get Hash Database
	hashModels := hashrecs.HashRecs{
		FilePath: HASH_FILE_PATH,
		Name:     HASH_DB,
	}

	hashModels.Load()

	//Syncing
	var responseModel []BackToCDSVendorResponse

	for i := 0; i < len(fileNames); i++ {
		//Get Json data from the file
		fileName := fileNames[i]
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileName)
		jsonString := string(jsonData)
		logFileName := "sync-" + fileName + ".log"

		// Unmarshal JSON to struct
		var vendors []WSVendor

		if err := json.Unmarshal([]byte(jsonData), &vendors); err != nil {
			message := "Failed[2]: error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.FAILURE, message, jsonString)
		}

		for j := 0; j < len(vendors); j++ {
			vendors[j].WeighbridgeSupplierID = PREFIX + vendors[j].WeighbridgeSupplierID
		}

		for j := 0; j < len(vendors); j++ {
			key := vendors[j].WeighbridgeSupplierID
			modelStr, _ := data_parser.ParseModelToString(vendors[j])
			hash := hashrecs.Hash(modelStr)
			preHash := hashModels.GetHash(key)

			if preHash == "" {
				utils.Console(fmt.Sprintf("Insert vendor: %s", key))

				isSuccess, err, result := InsertToNav(vendors[j])
				reqPayload, _ := data_parser.ParseJsonToXml(vendors[j])

				if err != nil {
					message := err.Error()
					utils.Console(message)
					payloads := fmt.Sprintf(`Request:\n %s`, reqPayload)
					logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.FAILURE, message, payloads)
				}

				if isSuccess {
					// Convert the string to a byte slice
					xmlData := []byte(result.(string))

					// Map Go struct to XML
					var parseModel CreateResultVendor
					err = xml.Unmarshal(xmlData, &parseModel)

					if err != nil {
						message := "Failed[6]: " + err.Error()
						utils.Console(message)
						payloads := fmt.Sprintf(`Request:\n %s \n\n Response:\n %s`, reqPayload, xmlData)
						logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.FAILURE, message, payloads)
					}

					// append success hash map and save hash map
					// Update the Hash field for a specific key
					hashModels.Set(key, hashrecs.HashRec{
						Hash:  hash,
						NavID: parseModel.Body.CreateResult.WSVendor.No,
					})

					//Add successed to an array
					responseModel = append(responseModel, BackToCDSVendorResponse{
						VendorNo:              parseModel.Body.CreateResult.WSVendor.No,
						WeighbridgeSupplierID: parseModel.Body.CreateResult.WSVendor.WeighbridgeSupplierID,
					})
				}
			}

			if preHash != "" && preHash != hash {
				utils.Console(fmt.Sprintf("Update vendor: %s", key))
				// @TODO: Update the vendor
			}

			if preHash != "" && preHash == hash {
				utils.Console(fmt.Sprintf("Skip vendor: %s", key))
			}
		}

		//Move to done file
		err = filesystem.MoveFile(fileName, PENDING_FILE_PATH, DONE_FILE_PATH)

		if err != nil {
			message := "Failed[5]: " + err.Error()
			utils.Console(message)
			logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.FAILURE, message, "")
		} else {
			//isSuccessfullySavedToFile = true
			message := "File successfully moved to the proccessed folder"
			utils.Console(message)
			logger.AddToLog(VENDOR_LOG_PATH, logFileName, logger.SUCCESS, message, fileName)
		}
	}

	//Save to Hash Folder
	hashModels.Save()

	//Bulk save
	//After syncing all files then send response back to CDS
	//utils.Console("responseModel------> ", len(responseModel))

	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func ReSync() {
	//Path
	HASH_FILE_PATH := utils.VENDOR_HASH_FILE_PATH
	HASH_DB := utils.VENDOR_HASH_DB

	//Get Hash Database
	hashModels := hashrecs.HashRecs{
		FilePath: HASH_FILE_PATH,
		Name:     HASH_DB,
	}
	hashModels.Load()

	//Mapping Model
	var responseModel []BackToCDSVendorResponse
	for key, value := range hashModels.Recs {
		responseModel = append(responseModel, BackToCDSVendorResponse{
			VendorNo:              key,
			WeighbridgeSupplierID: value.NavID,
		})
	}

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
