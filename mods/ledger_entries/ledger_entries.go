package ledgerentries

import (
	"encoding/json"
	"fmt"
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
	FETCH_URL := config.Config.LedgerEntries.Fetch.URL
	TOKEN_KEY := config.Config.Invoice.Fetch.APIKey
	PENDING_FILE_PATH := utils.LEDGER_ENTRIES_PENDING_FILE_PATH
	PENDING_LOG_FILE_PATH := utils.LEDGER_ENTRIES_DONE_LOG_FILE_PATH
	PENDING_FAILURE := utils.LEDGER_ENTRIES_PENDING_FAILURE
	PENDING_SUCCESS := utils.LEDGER_ENTRIES_PENDING_SUCCESS

	//Fetch LEDGER_ENTRIES data
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
		message := "Fetch: Successfully saved Ledger Entries to pending file"
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_SUCCESS, timestamp+".json", message, "")
	}
}

func Sync3() {
	//Path
	PENDING_FILE_PATH := utils.LEDGER_ENTRIES_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.LEDGER_ENTRIES_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.LEDGER_ENTRIES_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.LEDGER_ENTRIES_DONE_FAILURE
	DONE_SUCCESS := utils.LEDGER_ENTRIES_DONE_SUCCESS
	HASH_FILE_PATH := utils.LEDGER_ENTRIES_HASH_FILE_PATH
	HASH_DB := utils.LEDGER_ENTRIES_HASH_DB

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
	var responseModel []BackToCDSLedgerEntriesResponse
	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])
		if err != nil {
			message := "Failed:Sync:1 Could not read file -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
		}
		jsonString := string(jsonData)

		// Unmarshal JSON to struct
		var ledger_entries []LedgerEntriesCreate
		if err := json.Unmarshal([]byte(jsonData), &ledger_entries); err != nil {
			message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
		}

		for j := 0; j < len(ledger_entries); j++ {
			if ledger_entries[j].VendorPayment.AppliesToDocNo != nil && ledger_entries[j].VendorPayment.AccountNo != nil {
				key := *ledger_entries[j].VendorPayment.AppliesToDocNo
				vendorNo := *ledger_entries[j].VendorPayment.AccountNo
				paymentId := ledger_entries[j].PaymentID
				modelStr, _ := data_parser.ParseModelToString(ledger_entries[j])
				hash := hashrecs.Hash(modelStr)
				preHash := hashModels.GetHash(key)

				if preHash == "" {
					isSuccessCreation, err, resultCreate := InsertToNav(ledger_entries[j])
					if err != nil {
						message := "Failed:Sync:3 -> " + err.Error()
						utils.Console(message)
						logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
					}

					if isSuccessCreation {
						// Map Go struct to XML
						createLedgerEntryRes, err := UnmarshelCreateLedgerEntryResponse(resultCreate)
						if err != nil {
							message := "Failed:Sync:4 Error unmarshaling JSON -> " + err.Error()
							utils.Console(message)
							logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
						}

						isSuccessPost, err, resultPost := PostLedgerEntriesAfterCreation(createLedgerEntryRes)
						if err != nil {
							message := "Failed:Sync:5 -> " + err.Error()
							utils.Console(message)
							logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, jsonString)
						}

						if isSuccessPost {
							postLedgerEntryRes, err := UnmarshelCreateLedgerEntryResponse(resultPost)
							if err != nil {
								message := "Failed:Sync:6 " + err.Error()
								utils.Console(message)
							}

							//map
							// vendorNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.BuyFromVendorNo
							// purchaseInvoiceNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No
							documentNo := postLedgerEntryRes.Body.CreateResult.VendorPayment.DocumentNo

							// append success hash map and save hash map
							// Update the Hash field for a specific key
							hashModels.Set(key, hashrecs.HashRec{
								Hash:       hash,
								NavID:      vendorNo,
								DocumentNo: fmt.Sprintf("%i", documentNo),
								PaymentId:  paymentId,
							})
							if err != nil {
								utils.Console("UpdateHashInModel::Error:", err)
							}
							utils.Console("hashMaps", hashModels)

							//Add successed to an array
							responseModel = append(responseModel, BackToCDSLedgerEntriesResponse{
								//RefundId:          refundId,
								//PurchaseInvoiceNo: purchaseInvoiceNo,
								VendorNo:   vendorNo,
								DocumentNo: fmt.Sprintf("%i", documentNo),
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
			message := "Failed:Sync:4 " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
		} else {
			message := "Sync: File moved successfully"
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
	HASH_FILE_PATH := utils.LEDGER_ENTRIES_HASH_FILE_PATH
	HASH_DB := utils.LEDGER_ENTRIES_HASH_DB

	//Get Hash Database
	hashModels := hashrecs.HashRecs{
		FilePath: HASH_FILE_PATH,
		Name:     HASH_DB,
	}
	hashModels.Load()

	//Mapping Model
	var responseModel []BackToCDSLedgerEntriesResponse
	for _, value := range hashModels.Recs {
		responseModel = append(responseModel, BackToCDSLedgerEntriesResponse{
			//RefundId:          refundId,
			//PurchaseInvoiceNo: purchaseInvoiceNo,
			VendorNo:   value.NavID,
			DocumentNo: value.DocumentNo,
		})
	}

	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func sendToCDS(responseModel []BackToCDSLedgerEntriesResponse) {
	//Path
	RESPONSE_URL := config.Config.LedgerEntries.Save.URL
	// TOKEN_KEY := config.Config.Invoice.Fetch.APIKey

	// //Save Response vendor data to CDS
	// response, err := manager.Fetch(RESPONSE_URL, normalapi.POST, TOKEN_KEY, responseModel)
	// if err != nil {
	// 	message := "Failed:Fetch:1 " + err.Error()
	// 	utils.Console(message)
	// 	//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	// }
	// utils.Console(response)

	utils.Console("Successfully send to CDS system from nav ---> ledger entry: ", RESPONSE_URL)
	utils.Console(responseModel)
}
