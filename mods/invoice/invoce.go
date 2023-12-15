package invoice

import (
	"encoding/json"
	"nav_sync/config"
	"nav_sync/logger"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/mods/ahelpers/manager"
	normalapi "nav_sync/mods/ahelpers/normal_api"
	data_parser "nav_sync/mods/ahelpers/parser"
	"nav_sync/mods/hashrecs"
	"nav_sync/utils"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.Invoice.Fetch.URL
	TOKEN_KEY := config.Config.Invoice.Fetch.APIKey
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	PENDING_LOG_FILE_PATH := utils.INVOICE_PENDING_LOG_FILE_PATH
	PENDING_FAILURE := utils.INVOICE_PENDING_FAILURE
	PENDING_SUCCESS := utils.INVOICE_PENDING_SUCCESS

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
		message := "Fetch: Successfully saved invoice to pending file"
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_SUCCESS, timestamp+".json", message, "")
	}

}

func Sync3() {
	//Path
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.INVOICE_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.INVOICE_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS
	HASH_FILE_PATH := utils.INVOICE_HASH_FILE_PATH
	HASH_DB := utils.INVOICE_HASH_DB
	FAKE_PREFIX := config.Config.Invoice.FakePrefix
	FAKE_INSERT := config.Config.Invoice.FakeInsert

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(PENDING_FILE_PATH)
	if err != nil {
		message := "Failed:Sync:1 " + err.Error()
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, "", message, "")
	}

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
	var responseModel []BackToCDSInvoiceResponse
	for i := 0; i < len(fileNames); i++ {
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])
		if err != nil {
			message := "Failed:Sync:1 Could not read file -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
		}
		jsonString := string(jsonData)

		// Unmarshal JSON to struct
		var invoices []WSPurchaseInvoicePage
		if err := json.Unmarshal([]byte(jsonData), &invoices); err != nil {
			message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
		}

		if FAKE_INSERT {
			for j := 0; j < len(invoices); j++ {
				// 250500001
				invoices[j].VendorInvoiceNo = FAKE_PREFIX + invoices[j].VendorInvoiceNo
			}
		}

		for j := 0; j < len(invoices); j++ {
			if invoices[j].BuyFromVendorNo != nil {
				VendorInvoiceNoStr := invoices[j].VendorInvoiceNo //strconv.Itoa(invoices[j].VendorInvoiceNo)
				key := VendorInvoiceNoStr
				refundId := invoices[j].RefundId
				modelStr, _ := data_parser.ParseModelToString(invoices[j])
				hash := hashrecs.Hash(modelStr)
				preHash := hashModels.GetHash(key)

				if preHash == "" {
					isSuccessCreation, err, resultCreate := InsertToNav(invoices[j])
					if err != nil {
						message := "Failed:Sync:3 -> " + err.Error()
						utils.Console(message)
						logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
					}

					if isSuccessCreation {
						// Map Go struct to XML
						createInvoiceRes, err := UnmarshelCreateInvoiceResponse(resultCreate)
						if err != nil {
							message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
							utils.Console(message)
							logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
						}
						isSuccessPost, err, resultPost := PostToNavAfterInsert(createInvoiceRes)
						if err != nil {
							message := "Failed:Sync:4 -> " + err.Error()
							utils.Console(message)
							logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
						}

						if isSuccessPost {
							postInvoiceRes, err := UnmarshelPostInvoiceResponse(resultPost)
							if err != nil {
								message := "Failed:Sync:6 " + err.Error()
								utils.Console(message)
							}

							//map
							vendorNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.BuyFromVendorNo
							purchaseInvoiceNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No
							documentNo := postInvoiceRes.Body.ReturnValue

							// append success hash map and save hash map
							// Update the Hash field for a specific key
							hashModels.Set(key, hashrecs.HashRec{
								Hash:       hash,
								NavID:      vendorNo,
								InvoiceNo:  purchaseInvoiceNo,
								DocumentNo: documentNo,
								RefundId:   refundId,
							})
							if err != nil {
								utils.Console("UpdateHashInModel::Error:", err)
							}
							utils.Console("hashMaps", hashModels)

							responseModel = append(responseModel, BackToCDSInvoiceResponse{
								RefundId:          refundId,
								VendorNo:          vendorNo,
								PurchaseInvoiceNo: purchaseInvoiceNo,
								DocumentNo:        documentNo,
							})
						}
					}
				}

				if preHash != "" && preHash != hash {
					// @TODO: Update the vendor
				}
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
	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func ReSync() {
	//Path
	HASH_FILE_PATH := utils.INVOICE_HASH_FILE_PATH
	HASH_DB := utils.INVOICE_HASH_DB

	//Get Hash Database
	hashModels := hashrecs.HashRecs{
		FilePath: HASH_FILE_PATH,
		Name:     HASH_DB,
	}
	hashModels.Load()

	//Mapping Model
	var responseModel []BackToCDSInvoiceResponse
	for _, value := range hashModels.Recs {
		responseModel = append(responseModel, BackToCDSInvoiceResponse{
			RefundId:          value.RefundId,
			VendorNo:          value.NavID,
			PurchaseInvoiceNo: value.InvoiceNo,
			DocumentNo:        value.DocumentNo,
		})
	}

	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func sendToCDS(responseModel []BackToCDSInvoiceResponse) {
	//Path
	RESPONSE_URL := config.Config.Invoice.Save.URL
	TOKEN_KEY := config.Config.Invoice.Fetch.APIKey

	//Save Response vendor data to CDS
	response, err := manager.Fetch(RESPONSE_URL, normalapi.POST, TOKEN_KEY, responseModel)
	if err != nil {
		message := "Failed:sendToCDS:1 " + err.Error()
		utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	}
	utils.Console(response)

	// utils.Console("Successfully send to CDS system from nav ---> invoice: ", RESPONSE_URL)
	// utils.Console(responseModel)
}
