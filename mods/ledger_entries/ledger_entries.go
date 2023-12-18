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
	"strconv"
	"strings"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.LedgerEntries.Fetch.URL
	TOKEN_KEY := config.Config.Invoice.Fetch.APIKey
	IS_EMPTY := config.Config.LedgerEntries.EmptyLogs
	PENDING_FILE_PATH := utils.LEDGER_ENTRIES_PENDING_FILE_PATH
	LOG_PATH := utils.LEDGER_ENTRIES_LOG_PATH
	EMPTY_LOG_PATH := utils.LEDGER_ENTRIES_EMPTY_LOG_PATH
	EMPTY_LOG_DB := utils.EMPTY_LOG_DB

	utils.Console("Start fetching payments")

	//get current timestamp
	timestamp := utils.GetCurrentTime()
	logFileName := "fetch-" + timestamp + ".log"

	//Fetch payment data
	response, err := manager.Fetch(FETCH_URL, normalapi.GET, TOKEN_KEY, nil)
	if err != nil {
		message := "Failed[1]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, FETCH_URL)
		return
	}

	//Checking if cointains data
	var ledgers []LedgerEntriesCreate
	ledgers, _ = UnmarshalStringToLedgerEntries(response)
	if IS_EMPTY && len(ledgers) < 1 {
		//Save logs
		logger.AddToLog(EMPTY_LOG_PATH, EMPTY_LOG_DB+".log", logger.EMPTY, "Fetched ledgers with empty", "")
		return
	}

	//Save to pending file
	err = filesystem.Save(PENDING_FILE_PATH, timestamp, response)

	if err != nil {
		message := "Failed[2]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, "")
	} else {
		message := "fetched payments and saved to a file"
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileName, logger.SUCCESS, message, PENDING_FILE_PATH+timestamp+".json")
	}
}

func Sync3() {
	//Path
	PENDING_FILE_PATH := utils.LEDGER_ENTRIES_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.LEDGER_ENTRIES_DONE_FILE_PATH
	LOG_PATH := utils.LEDGER_ENTRIES_LOG_PATH
	HASH_FILE_PATH := utils.LEDGER_ENTRIES_HASH_FILE_PATH
	HASH_DB := utils.LEDGER_ENTRIES_HASH_DB

	utils.Console("Start sending invoices to NAV")

	timestamp := utils.GetCurrentTime()
	logFileGeneral := "sync-pre-" + timestamp + ".log"

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(PENDING_FILE_PATH)

	if err != nil {
		message := "Failed[1]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileGeneral, logger.FAILURE, message, "")
		return
	}

	if fileNames == nil || len(fileNames) < 1 {
		message := "No pending files found"
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileGeneral, logger.SUCCESS, "Skipped", message)
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
		fileName := fileNames[i]
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileName)
		logFileName := "sync-" + strings.Replace(fileName, ".json", "", 1) + ".log"

		if err != nil {
			message := "Failed[1]: could not read file -> " + err.Error()
			utils.Console(message)
			logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, PENDING_FILE_PATH+fileName)
			continue
		}

		// Unmarshal JSON to struct
		var ledger_entries []LedgerEntriesCreate

		if err := json.Unmarshal([]byte(jsonData), &ledger_entries); err != nil {
			message := "Failed[2]: error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, DONE_FILE_PATH+fileName)
			return
		}

		for j := 0; j < len(ledger_entries); j++ {
			if ledger_entries[j].VendorPayment.AppliesToDocNo != nil && ledger_entries[j].VendorPayment.AccountNo != nil {
				paymentId := ledger_entries[j].PaymentID
				// key := *ledger_entries[j].VendorPayment.AppliesToDocNo
				key := strconv.Itoa(paymentId)
				vendorNo := *ledger_entries[j].VendorPayment.AccountNo
				modelStr, _ := data_parser.ParseModelToString(ledger_entries[j])
				hash := hashrecs.Hash(modelStr)
				preHash := hashModels.GetHash(key)

				if preHash == "" {
					utils.Console(fmt.Sprintf("Insert payment: %s", key))

					isSuccessCreation, err, resultCreate := InsertToNav(ledger_entries[j])
					reqPayload, _ := data_parser.ParseJsonToXml(ledger_entries[j])

					if err != nil {
						message := err.Error()
						utils.Console(message)
						payloads := fmt.Sprintf(`Request:\n %s`, reqPayload)
						logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, payloads)
						continue
					}

					if isSuccessCreation {
						// Map Go struct to XML
						createLedgerEntryRes, err := UnmarshelCreateLedgerEntryResponse(resultCreate)

						if err != nil {
							message := "Failed[2]: Error unmarshaling JSON -> " + err.Error()
							utils.Console(message)
							logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, fileName)
							continue
						}

						// utils.Console(fmt.Sprintf("Post invoice: %s", createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No))

						isSuccessPost, err, _ := PostLedgerEntriesAfterCreation(createLedgerEntryRes)

						if err != nil {
							message := "Failed[4]: " + err.Error()
							utils.Console(message)
							logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, fileName)
							continue
						}

						if isSuccessPost {
							//map
							// vendorNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.BuyFromVendorNo
							// purchaseInvoiceNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No
							documentNo := createLedgerEntryRes.Body.CreateResult.VendorPayment.DocumentNo
							documentNoStr := strconv.Itoa(documentNo)

							// append success hash map and save hash map
							// Update the Hash field for a specific key
							hashModels.Set(key, hashrecs.HashRec{
								Hash:       hash,
								NavID:      vendorNo,
								DocumentNo: documentNoStr,
								PaymentId:  paymentId,
							})

							//Add successed to an array
							responseModel = append(responseModel, BackToCDSLedgerEntriesResponse{
								//RefundId:          refundId,
								//PurchaseInvoiceNo: purchaseInvoiceNo,
								PaymentId:  paymentId,
								VendorNo:   vendorNo,
								DocumentNo: documentNoStr,
							})

						}

					}
				}

				if preHash != "" && preHash != hash {
					utils.Console(fmt.Sprintf("Update payment: %s", key))
					// @TODO: Update the payment
				}

				if preHash != "" && preHash == hash {
					utils.Console(fmt.Sprintf("Skip payment: %s", key))
				}
			}
		}

		//Move to done file
		err = filesystem.MoveFile(fileName, PENDING_FILE_PATH, DONE_FILE_PATH)

		if err != nil {
			message := "Failed[5]: " + err.Error()
			utils.Console(message)
			logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, "")
		} else {
			//isSuccessfullySavedToFile = true
			message := "File successfully moved to the proccessed folder"
			utils.Console(message)
			logger.AddToLog(LOG_PATH, logFileName, logger.SUCCESS, message, fileName)
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
			PaymentId:  value.PaymentId,
		})
	}

	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func sendToCDS(responseModel []BackToCDSLedgerEntriesResponse) {
	//Path
	RESPONSE_URL := config.Config.LedgerEntries.Save.URL
	TOKEN_KEY := config.Config.Invoice.Fetch.APIKey

	//Save Response vendor data to CDS
	_, err := manager.Fetch(RESPONSE_URL, normalapi.POST, TOKEN_KEY, responseModel)

	if err != nil {
		message := "Failed:Fetch:1 " + err.Error()
		utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	}

	// utils.Console("Successfully send to CDS system from nav ---> ledger entry: ", RESPONSE_URL)
	// utils.Console(responseModel)
}
