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
	"nav_sync/utils"
)

func Fetch() {
	//Path
	FETCH_URL := config.Config.Invoice.Fetch.URL
	PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	PENDING_LOG_FILE_PATH := utils.INVOICE_PENDING_LOG_FILE_PATH
	PENDING_FAILURE := utils.INVOICE_PENDING_FAILURE
	PENDING_SUCCESS := utils.INVOICE_PENDING_SUCCESS

	//Fetch vendor data
	response, err := manager.Fetch(FETCH_URL, normalapi.GET)
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

// func Sync() {
// 	//Eg.
// 	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
// 	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
// 	url := config.Config.Invoice.Sync.URL
// 	xmlPayload := `
// 		<Envelope
// 			xmlns="http://schemas.xmlsoap.org/soap/envelope/">
// 			<Body>
// 				<Create
// 					xmlns="urn:microsoft-dynamics-schemas/page/wspurchaseinvoicepage">
// 					<WSPurchaseInvoicePage>
// 						<Buy_from_Vendor_No>TEST2</Buy_from_Vendor_No>
// 						<Vendor_Invoice_No>124A</Vendor_Invoice_No>
// 						<Buy_from_Vendor_Name>intuji 2</Buy_from_Vendor_Name>
// 						<!-- Optional -->
// 						<PurchLines>
// 							<Purch_Invoice_Line>
// 								<Type>Item</Type>
// 								<No>250MS-12MTR-006</No>
// 								<Quantity>1</Quantity>
// 								<Unit_Price_LCY>10</Unit_Price_LCY>
// 								<Location_Code>BI</Location_Code>
// 							</Purch_Invoice_Line>
// 						</PurchLines>
// 					</WSPurchaseInvoicePage>
// 				</Create>
// 			</Body>
// 		</Envelope>
// 	`

// 	utils.Console(url)
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

	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(PENDING_FILE_PATH, fileNames[i])

		//Insert invoice
		response, err := insertInvoice(jsonData)
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
			match := utils.MatchRegexExpression(resultStr, `<WSPurchaseInvoicePage[^>]*>`)
			matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

			// Print the result
			if !match && matchFault {
				message := fmt.Sprintf("Failed:Sync:4 XML string does not contain <WSPurchaseInvoicePage> element: ", resultStr)
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
			} else {
				isSuccessCreation = true
			}
		}

		isSuccessPost := false
		if isSuccessCreation {
			response, err = postInvoiceAfterCreation(response)
			if err != nil {
				isSuccessPost = false
				message := "Failed:Sync:5 " + err.Error()
				utils.Console(message)
				logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, "")
			} else {
				utils.Console(response)
				resultStr, ok := response.(string)
				if !ok {
					// The type assertion failed
					message := fmt.Sprintf("Failed:Sync:6 Could not convert to string: ", response)
					utils.Console(message)
					logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
				}
				match := utils.MatchRegexExpression(resultStr, `<PostPurchaseResult[^>]*>`)
				matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

				// Print the result
				if !match && matchFault {
					message := fmt.Sprintf("Failed:Sync:7 XML string does not contain <PostPurchaseResult> element: ", resultStr)
					utils.Console(message)
					logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileNames[i], message, resultStr)
				} else {
					isSuccessPost = true
				}
			}
		}

		//Move to done file
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

	// Type assertion to get the underlying string
	str, ok := stringData.(string)
	if !ok {
		return nil, errors.New("postInvoiceAfterCreation: Conversion failed: not a string")
	}

	// Convert the string to a byte slice
	xmlData := []byte(str)

	// Map Go struct to XML
	var envelope PostInvoiceEnvelope
	err := xml.Unmarshal(xmlData, &envelope)
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
