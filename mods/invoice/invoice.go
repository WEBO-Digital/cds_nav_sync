package invoice

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
	"strings"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.Invoice.Fetch.URL
	TOKEN_KEY := config.Config.Invoice.Fetch.APIKey
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	LOG_PATH := utils.INVOICE_LOG_PATH

	utils.Console("Start fetching invoices")

	//get current timestamp
	timestamp := utils.GetCurrentTime()
	logFileName := "fetch-" + timestamp + ".log"

	//Fetch invoice data
	response, err := manager.Fetch(FETCH_URL, normalapi.GET, TOKEN_KEY, nil)

	if err != nil {
		message := "Failed[1]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, FETCH_URL)
		return
	}

	//Save to pending file
	err = filesystem.Save(PENDING_FILE_PATH, timestamp, response)

	if err != nil {
		message := "Failed[2]: " + err.Error()
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, "")
	} else {
		message := "fetched invoices and saved to a file"
		utils.Console(message)
		logger.AddToLog(LOG_PATH, logFileName, logger.SUCCESS, message, PENDING_FILE_PATH+timestamp+".json")
	}
}

func Sync3() {
	//Path
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.INVOICE_DONE_FILE_PATH
	LOG_PATH := utils.INVOICE_LOG_PATH
	HASH_FILE_PATH := utils.INVOICE_HASH_FILE_PATH
	HASH_DB := utils.INVOICE_HASH_DB
	PREFIX := config.Config.Invoice.Prefix

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
	var responseModel []BackToCDSInvoiceResponse

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

		// jsonString := string(jsonData)

		// Unmarshal JSON to struct
		var invoices []WSPurchaseInvoicePage

		if err := json.Unmarshal([]byte(jsonData), &invoices); err != nil {
			message := "Failed[2]: error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, DONE_FILE_PATH+fileName)
			return
		}

		for j := 0; j < len(invoices); j++ {
			invoices[j].VendorInvoiceNo = PREFIX + invoices[j].VendorInvoiceNo
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
					utils.Console(fmt.Sprintf("Insert invoice: %s", key))

					isSuccessCreation, err, resultCreate := InsertToNav(invoices[j])
					reqPayload, _ := data_parser.ParseJsonToXml(invoices[j])

					if err != nil {
						message := err.Error()
						utils.Console(message)
						payloads := fmt.Sprintf(`Request:\n %s`, reqPayload)
						logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, payloads)
						continue
					}

					if isSuccessCreation {
						// Map Go struct to XML
						createInvoiceRes, err := UnmarshelCreateInvoiceResponse(resultCreate)

						if err != nil {
							message := "Failed[2]: Error unmarshaling JSON -> " + err.Error()
							utils.Console(message)
							logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, fileName)
							continue
						}

						utils.Console(fmt.Sprintf("Post invoice: %s", createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No))

						isSuccessPost, err, resultPost := PostToNavAfterInsert(createInvoiceRes)

						if err != nil {
							message := "Failed[4]: " + err.Error()
							utils.Console(message)
							logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, fileName)
							continue
						}

						if isSuccessPost {
							postInvoiceRes, err := UnmarshelPostInvoiceResponse(resultPost)

							if err != nil {
								message := "Failed[6]: " + err.Error()
								utils.Console(message)
								logger.AddToLog(LOG_PATH, logFileName, logger.FAILURE, message, fileName)
								continue
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
					utils.Console(fmt.Sprintf("Update invoice: %s", key))
					// @TODO: Update the invoice
				}

				if preHash != "" && preHash == hash {
					utils.Console(fmt.Sprintf("Skip invoice: %s", key))
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
	_, err := manager.Fetch(RESPONSE_URL, normalapi.POST, TOKEN_KEY, responseModel)
	if err != nil {
		message := "Failed:sendToCDS:1 " + err.Error()
		utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	}

	// utils.Console("Successfully send to CDS system from nav ---> invoice: ", RESPONSE_URL)
	// utils.Console(responseModel)
}
