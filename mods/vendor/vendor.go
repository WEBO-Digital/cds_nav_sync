package vendor

import (
	"encoding/json"
	"encoding/xml"
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

func Sync3() {
	//Path
	PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.VENDOR_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.VENDOR_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS
	HASH_FILE_PATH := utils.VENDOR_HASH_FILE_PATH
	HASH_DB := utils.VENDOR_HASH_DB
	FAKE_PREFIX := config.Config.Vendor.FakePrefix
	FAKE_INSERT := config.Config.Vendor.FakeInsert

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
	hashModels := hashrecs.HashRecs{
		FilePath: HASH_FILE_PATH,
		Name:     HASH_DB,
	}
	hashModels.Load()

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

		if FAKE_INSERT {
			for j := 0; j < len(vendors); j++ {
				// 250500001
				vendors[j].WeighbridgeSupplierID = FAKE_PREFIX + vendors[j].WeighbridgeSupplierID
			}
		}

		for j := 0; j < len(vendors); j++ {
			key := vendors[j].WeighbridgeSupplierID
			modelStr, _ := data_parser.ParseModelToString(vendors[j])
			hash := hashrecs.Hash(modelStr)
			preHash := hashModels.GetHash(key)

			if preHash == "" {
				isSuccess, err, result := InsertToNav(vendors[j])

				if err != nil {
					message := err.Error()
					utils.Console(message)
					logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
				}

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

					// append success hash map and save hash map
					// Update the Hash field for a specific key
					hashModels.Set(key, hashrecs.HashRec{
						Hash:  hash,
						NavID: parseModel.Body.CreateResult.WSVendor.No,
					})
					if err != nil {
						utils.Console("UpdateHashInModel::Error:", err)
					}
					utils.Console("hashMaps", hashModels)

					//Add successed to an array
					responseModel = append(responseModel, BackToCDSVendorResponse{
						VendorNo:              parseModel.Body.CreateResult.WSVendor.No,
						WeighbridgeSupplierID: parseModel.Body.CreateResult.WSVendor.WeighbridgeSupplierID,
					})
				}
			}

			if preHash != "" && preHash != hash {
				// @TODO: Update the vendor
			}
		}

		//Move to done file
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
	}

	//Save to Hash Folder
	hashModels.Save()

	//Bulk save
	//After syncing all files then send response back to CDS
	utils.Console("responseModel------> ", len(responseModel))

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
