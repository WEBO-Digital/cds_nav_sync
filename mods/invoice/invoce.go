package invoice

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"nav_sync/config"
	"nav_sync/logger"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/mods/ahelpers/manager"
	navapi "nav_sync/mods/ahelpers/nav_api"
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

func Sync() {
	//Path
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.INVOICE_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.INVOICE_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(PENDING_FILE_PATH)
	if err != nil {
		message := "Failed:Sync:1 " + err.Error()
		utils.Console(message)
		logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, "", message, "")
	}

	//mods.Console(fileNames)

	if fileNames == nil || len(fileNames) < 1 {
		return
	}

	var responseModel []BackToCDSInvoiceResponse
	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])

		//Insert invoice
		responseCreate, err := insertInvoice(jsonData)
		isSuccessCreation := false
		if err != nil {
			isSuccessCreation = false
			message := "Failed:Sync:2 " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
		} else {
			utils.Console(responseCreate)
			resultCreateStr, ok := responseCreate.(string)
			if !ok {
				// The type assertion failed
				message := fmt.Sprintf("Failed:Sync:3 Could not convert to string: ", responseCreate)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultCreateStr)
			}
			match := utils.MatchRegexExpression(resultCreateStr, `<WSPurchaseInvoicePage[^>]*>`)
			matchFault := utils.MatchRegexExpression(resultCreateStr, `<faultcode[^>]*>`)

			// Print the result
			if !match && matchFault {
				message := fmt.Sprintf("Failed:Sync:4 XML string does not contain <WSPurchaseInvoicePage> element: ", resultCreateStr)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultCreateStr)
			} else {
				isSuccessCreation = true
			}
		}

		isSuccessPost := false
		var responsePost interface{}
		if isSuccessCreation {
			responsePost, err = postInvoiceAfterCreation(responseCreate)
			if err != nil {
				isSuccessPost = false
				message := "Failed:Sync:5 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
			} else {
				utils.Console(responsePost)
				resultPostStr, ok := responsePost.(string)
				if !ok {
					// The type assertion failed
					message := fmt.Sprintf("Failed:Sync:6 Could not convert to string: ", resultPostStr)
					utils.Console(message)
					logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultPostStr)
				}
				match := utils.MatchRegexExpression(resultPostStr, `<PostPurchaseResult[^>]*>`)
				matchFault := utils.MatchRegexExpression(resultPostStr, `<faultcode[^>]*>`)

				// Print the result
				if !match && matchFault {
					message := fmt.Sprintf("Failed:Sync:7 XML string does not contain <PostPurchaseResult> element: ", resultPostStr)
					utils.Console(message)
					logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultPostStr)
				} else {
					isSuccessPost = true
				}
			}
		}

		//Move to done file
		isSuccessfullySavedToFile := false
		if isSuccessCreation && isSuccessPost {
			err = filesystem.MoveFile(fileNames[i], PENDING_FILE_PATH, DONE_FILE_PATH)
			if err != nil {
				message := "Failed:Sync:4 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
			} else {
				isSuccessfullySavedToFile = true
				message := "Sync: File moved successfully"
				utils.Console(message)
				logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_SUCCESS, fileNames[i], message, "")
			}
		}

		//Add successed to an array
		if isSuccessfullySavedToFile {
			createInvoiceRes, err := unmarshelCreateInvoiceResponse(responseCreate)
			if err != nil {
				message := "Failed:Sync:6 " + err.Error()
				utils.Console(message)
			}

			postInvoiceRes, err := unmarshelPostInvoiceResponse(responsePost)
			if err != nil {
				message := "Failed:Sync:6 " + err.Error()
				utils.Console(message)
			}

			responseModel = append(responseModel, BackToCDSInvoiceResponse{
				VendorNo:          createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.BuyFromVendorNo,
				PurchaseInvoiceNo: createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No,
				DocumentNo:        postInvoiceRes.Body.ReturnValue,
			})
		}
	}

	//Bulk save
	//After syncing all files then send response back to CDS
	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func insertInvoice(jsonData []byte) (interface{}, error) {
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Invoice.Sync.URL

	// Unmarshal JSON to struct
	var invoice WSPurchaseInvoicePage
	if err := json.Unmarshal([]byte(jsonData), &invoice); err != nil {
		return nil, errors.New("insertInvoice: Error unmarshaling JSON -> " + err.Error())
	}

	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(invoice)
	if err != nil {
		return nil, errors.New("insertInvoice: Error mapping to XML -> " + err.Error())
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<Create xmlns="urn:microsoft-dynamics-schemas/page/wspurchaseinvoicepage">
						%s
					</Create>
				</Body>
			</Envelope>
		`,
		string(xmlData),
	)

	//Return the result
	utils.Console(xmlPayload)
	utils.Console("username: ", NTLM_USERNAME)
	utils.Console("username: ", NTLM_PASSWORD)
	utils.Console("URL: ", url)

	//Sync to Nav
	result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		return nil, errors.New("insertInvoice: " + err.Error())
	}
	return result, nil
}

func postInvoiceAfterCreation(stringData interface{}) (interface{}, error) {
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Invoice.Post.URL

	// Map Go struct to XML
	envelope, err := unmarshelCreateInvoiceResponse(stringData)
	if err != nil {
		return nil, errors.New("postInvoiceAfterCreation: Error decoding XML: " + err.Error())
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<PostPurchaseInvoice xmlns="urn:microsoft-dynamics-schemas/codeunit/WSPurchaseInvoice">
						<docNo>%s</docNo>
						<lastDocumentType>%s</lastDocumentType>
					</PostPurchaseInvoice>
				</Body>
			</Envelope>
		`,
		envelope.Body.CreateResult.WSPurchaseInvoicePage.No,
		"2",
	)

	//Return the result
	utils.Console(xmlPayload)
	utils.Console("username: ", NTLM_USERNAME)
	utils.Console("username: ", NTLM_PASSWORD)
	utils.Console("URL: ", url)

	//Sync to Nav
	result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		return nil, errors.New("postInvoiceAfterCreation: " + err.Error())
	}
	return result, nil
}

func unmarshelCreateInvoiceResponse(stringData interface{}) (PostInvoiceEnvelope, error) {
	var envelope PostInvoiceEnvelope
	// Type assertion to get the underlying string
	str, ok := stringData.(string)
	if !ok {
		return envelope, errors.New("unmarshelCreateInvoiceResponse: Conversion failed: not a string")
	}

	// Convert the string to a byte slice
	xmlData := []byte(str)

	// Map Go struct to XML
	err := xml.Unmarshal(xmlData, &envelope)
	if err != nil {
		return envelope, errors.New("unmarshelCreateInvoiceResponse: Error decoding XML: " + err.Error())
	}
	return envelope, nil
}

func unmarshelPostInvoiceResponse(stringData interface{}) (PostResponseInvoiceEnvelope, error) {
	var envelope PostResponseInvoiceEnvelope
	// Type assertion to get the underlying string
	str, ok := stringData.(string)
	if !ok {
		return envelope, errors.New("unmarshelPostInvoiceResponse: Conversion failed: not a string")
	}

	// Convert the string to a byte slice
	xmlData := []byte(str)

	// Map Go struct to XML
	err := xml.Unmarshal(xmlData, &envelope)
	if err != nil {
		return envelope, errors.New("unmarshelPostInvoiceResponse: Error decoding XML: " + err.Error())
	}
	return envelope, nil
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

func Sync2() {
	//Path
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.INVOICE_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.INVOICE_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS

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
	hashModels := GetHash()

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

		//var isSuccessArr []bool
		for j := 0; j < len(invoices); j++ {
			hashMaps, isHashed := CompareWithHash(invoices[j], hashModels)
			if !isHashed {
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
						vendorStr, _ := data_parser.ParseModelToString(invoices[j])
						refundId := invoices[j].RefundId
						vendorInvoiceNo := fmt.Sprintf("%i", invoices[j].VendorInvoiceNo)
						vendorNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.BuyFromVendorNo
						purchaseInvoiceNo := createInvoiceRes.Body.CreateResult.WSPurchaseInvoicePage.No
						documentId := postInvoiceRes.Body.ReturnValue

						// append success hash map and save hash map
						// Update the Hash field for a specific key
						err = UpdateHashInModel(hashMaps, vendorInvoiceNo, utils.ComputeMD5(vendorStr), vendorNo, purchaseInvoiceNo, documentId, refundId)
						if err != nil {
							utils.Console("UpdateHashInModel::Error:", err)
						}
						utils.Console("hashMaps", hashMaps)

						responseModel = append(responseModel, BackToCDSInvoiceResponse{
							RefundId:          refundId,
							VendorNo:          vendorNo,
							PurchaseInvoiceNo: purchaseInvoiceNo,
							DocumentNo:        documentId,
						})
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
	}

	//Save to Hash Folder
	val, err := SaveHashLogs(hashModels)
	if err != nil {
		utils.Console("SaveHashLogs::err", err)
	}
	utils.Console("SaveHashLogs::val", val)

	//Bulk save
	//After syncing all files then send response back to CDS
	if len(responseModel) > 0 {
		sendToCDS(responseModel)
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

		// for j := 0; j < len(invoices); j++ {
		// 	// 250500001
		// 	invoices[j].VendorInvoiceNo = 26959000 + j + 1
		// }

		for j := 0; j < len(invoices); j++ {
			if invoices[j].BuyFromVendorNo != nil {
				key := fmt.Sprintf("%i", invoices[j].VendorInvoiceNo)
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
