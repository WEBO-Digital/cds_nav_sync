package ledgerentries

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
	"nav_sync/utils"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.LedgerEntries.Fetch.URL
	PENDING_FILE_PATH := utils.LEDGER_ENTRIES_PENDING_FILE_PATH
	PENDING_LOG_FILE_PATH := utils.LEDGER_ENTRIES_DONE_LOG_FILE_PATH
	PENDING_FAILURE := utils.LEDGER_ENTRIES_PENDING_FAILURE
	PENDING_SUCCESS := utils.LEDGER_ENTRIES_PENDING_SUCCESS

	//Fetch LEDGER_ENTRIES data
	response, err := manager.Fetch(FETCH_URL, normalapi.GET, nil)
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

func Sync() {
	//Path
	PENDING_FILE_PATH := utils.LEDGER_ENTRIES_PENDING_FILE_PATH
	DONE_FILE_PATH := utils.LEDGER_ENTRIES_DONE_FILE_PATH
	DONE_LOG_FILE_PATH := utils.LEDGER_ENTRIES_DONE_LOG_FILE_PATH
	DONE_FAILURE := utils.LEDGER_ENTRIES_DONE_FAILURE
	DONE_SUCCESS := utils.LEDGER_ENTRIES_DONE_SUCCESS

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

	var responseModel []BackToCDSLedgerEntriesResponse
	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])

		//Insert invoice
		response, err := insertLedgerEntries(jsonData)
		isSuccessCreation := false
		if err != nil {
			isSuccessCreation = false
			message := "Failed:Sync:2 " + err.Error()
			utils.Console(message)
			logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
		} else {
			utils.Console(response)
			resultStr, ok := response.(string)
			if !ok {
				// The type assertion failed
				message := fmt.Sprintf("Failed:Sync:3 Could not convert to string: ", response)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
			}
			match := utils.MatchRegexExpression(resultStr, `<VendorPayment[^>]*>`)
			matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

			// Print the result
			if !match && matchFault {
				message := fmt.Sprintf("Failed:Sync:4 XML string does not contain <VendorPayment> element: ", resultStr)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
			} else {
				isSuccessCreation = true
			}
		}

		isSuccessPost := false
		var responsePost interface{}
		if isSuccessCreation {
			responsePost, err = postLedgerEntriesAfterCreation(responsePost)
			if err != nil {
				isSuccessPost = false
				message := "Failed:Sync:5 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
			} else {
				utils.Console(responsePost)
				resultStr, ok := responsePost.(string)
				if !ok {
					// The type assertion failed
					message := fmt.Sprintf("Failed:Sync:6 Could not convert to string: ", responsePost)
					utils.Console(message)
					logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
				}
				match := utils.MatchRegexExpression(resultStr, `<PostPurchaseInvoice_Result[^>]*>`)
				matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

				// Print the result
				if !match && matchFault {
					message := fmt.Sprintf("Failed:Sync:7 XML string does not contain <PostPurchaseInvoice_Result> element: ", resultStr)
					utils.Console(message)
					logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
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
				message := "Sync: File moved successfully"
				utils.Console(message)
				logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_SUCCESS, fileNames[i], message, "")
			}
		}

		//Add successed to an array
		if isSuccessfullySavedToFile {
			postLedgerRes, err := unmarshelCreateLedgerEntryResponse(responsePost)
			if err != nil {
				message := "Failed:Sync:6 " + err.Error()
				utils.Console(message)
			}

			utils.Console(responsePost)
			responseModel = append(responseModel, BackToCDSLedgerEntriesResponse{
				DocumentNo: postLedgerRes.Body.CreateResult.VendorPayment.DocumentNo,
			})
		}
	}

	//Bulk save
	//After syncing all files then send response back to CDS
	if len(responseModel) > 0 {
		sendToCDS(responseModel)
	}
}

func insertLedgerEntries(jsonData []byte) (interface{}, error) {
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.LedgerEntries.Sync.URL

	// Unmarshal JSON to struct
	var ledger_entries_create LedgerEntriesCreate
	if err := json.Unmarshal([]byte(jsonData), &ledger_entries_create); err != nil {
		return nil, errors.New("insertLedgerEntries: Error unmarshaling JSON -> " + err.Error())
	}

	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(ledger_entries_create.VendorPayment)
	if err != nil {
		return nil, errors.New("insertLedgerEntries: Error mapping to XML -> " + err.Error())
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<Create xmlns="urn:microsoft-dynamics-schemas/page/vendorpayment">
						<CurrentJnlBatchName>%s</CurrentJnlBatchName>
						%s
					</Create>
				</Body>
			</Envelope>
		`,
		ledger_entries_create.CurrentJnlBatchName,
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
		return nil, errors.New("insertLedgerEntries: " + err.Error())
	}
	return result, nil
}

func postLedgerEntriesAfterCreation(stringData interface{}) (interface{}, error) {
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.LedgerEntries.Post.URL

	// Type assertion to get the underlying string
	// str, ok := stringData.(string)
	// if !ok {
	// 	return nil, errors.New("postLedgerEntriesAfterCreation: Conversion failed: not a string")
	// }

	// // Convert the string to a byte slice
	// xmlData := []byte(str)

	// // Map Go struct to XML
	// var envelope PostLedgerEntriesEnvelope
	// err := xml.Unmarshal(xmlData, &envelope)
	// if err != nil {
	// 	return nil, errors.New("postLedgerEntriesAfterCreation: Error decoding XML: " + err.Error())
	// }

	envelope, err := unmarshelCreateLedgerEntryResponse(stringData)
	if err != nil {
		message := "postLedgerEntriesAfterCreation: " + err.Error()
		utils.Console(message)
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<PostGenJNlLine xmlns="urn:microsoft-dynamics-schemas/codeunit/WSPurchaseInvoice">
						<docNo>%v</docNo>
					</PostGenJNlLine>
				</Body>
			</Envelope>
		`,
		envelope.Body.CreateResult.VendorPayment.DocumentNo,
	)

	//Return the result
	utils.Console(xmlPayload)
	utils.Console("username: ", NTLM_USERNAME)
	utils.Console("username: ", NTLM_PASSWORD)
	utils.Console("URL: ", url)

	//Sync to Nav
	result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		return nil, errors.New("postLedgerEntriesAfterCreation: " + err.Error())
	}
	return result, nil
}

func unmarshelCreateLedgerEntryResponse(stringData interface{}) (PostLedgerEntriesEnvelope, error) {
	var envelope PostLedgerEntriesEnvelope
	// Type assertion to get the underlying string
	str, ok := stringData.(string)
	if !ok {
		return envelope, errors.New("unmarshelCreateLedgerEntryResponse: Conversion failed: not a string")
	}

	// Convert the string to a byte slice
	xmlData := []byte(str)

	// Map Go struct to XML
	err := xml.Unmarshal(xmlData, &envelope)
	if err != nil {
		return envelope, errors.New("unmarshelCreateLedgerEntryResponse: Error decoding XML: " + err.Error())
	}
	return envelope, nil
}

func sendToCDS(responseModel []BackToCDSLedgerEntriesResponse) {
	//Path
	RESPONSE_URL := config.Config.Vendor.Save.URL

	// //Save Response vendor data to CDS
	// response, err := manager.Fetch(RESPONSE_URL, normalapi.POST, responseModel)
	// if err != nil {
	// 	message := "Failed:Fetch:1 " + err.Error()
	// 	utils.Console(message)
	// 	//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, "", message, "")
	// }
	// utils.Console(response)

	utils.Console("Successfully send to CDS system from nav ---> ledger entry: ", RESPONSE_URL)
	utils.Console(responseModel)
}
