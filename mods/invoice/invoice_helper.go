package invoice

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"nav_sync/config"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/mods/ahelpers/manager"
	navapi "nav_sync/mods/ahelpers/nav_api"
	data_parser "nav_sync/mods/ahelpers/parser"
	"nav_sync/utils"
)

func GetHash() HashInvoiceModel {
	HASH_FILE_PATH := utils.INVOICE_HASH_FILE_PATH
	golbalModel := make(HashInvoiceModel)

	//Get All the invoices pending data
	fileNames, err := filesystem.GetAllFiles(HASH_FILE_PATH)
	if err != nil {
		return golbalModel
	}
	utils.Console("************************GetHash****************************")
	utils.Console(fileNames)
	if fileNames == nil || len(fileNames) < 1 {
		return golbalModel //, nil
	}

	for i := 0; i < len(fileNames); i++ {
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(HASH_FILE_PATH, fileNames[i])
		if err != nil {
			return golbalModel
		}

		jsonString := string(jsonData)
		utils.Console("************************GetHash::inner****************************")
		utils.Console(jsonString)

		// Unmarshal JSON to struct
		localModel := make(HashInvoiceModel)
		if err := json.Unmarshal([]byte(jsonData), &localModel); err != nil {
			return golbalModel
		}

		// Append localModel to golbalModel
		for key, value := range localModel {
			golbalModel[key] = value
		}
	}

	return golbalModel
}

func CompareWithHash(invoice WSPurchaseInvoicePage, hash HashInvoiceModel) (HashInvoiceModel, bool) {
	//var hashRecord HashVendorModel
	hashRecord := hash
	vendorStr, err := data_parser.ParseModelToString(invoice)
	if err != nil {
		return hashRecord, true
	}
	hashStr := utils.ComputeMD5(vendorStr)
	for key, _ := range hash {
		//hash is already pushed to nav server
		if containsHash(hash, hashStr) {
			return hashRecord, true
		}

		if !containsHash(hash, hashStr) && key == *invoice.BuyFromVendorNo { //&& value.NavID == nil
			/** @TODO: UPDATE THE VEDOR **/
			return hashRecord, true
		}
	}

	//hash is not found then insert to nav server
	return hashRecord, false
}

func containsHash(model HashInvoiceModel, targetHash string) bool {
	for _, entry := range model {
		if entry.Hash == targetHash {
			return true
		}
	}
	return false
}

func UpdateHashInModel(hashModel HashInvoiceModel, key, hash string, navId string, purchaseInvoiceNo string, documentNo string, refundId int) error {
	// Initialize the map if it is nil
	if hashModel == nil {
		hashModel = make(HashInvoiceModel)
	}

	// Check if the key exists in the map
	if entry, exists := hashModel[key]; exists {
		// Update the Hash field with the new value
		entry.Hash = hash
		entry.NavID = navId
		entry.InvoiceNo = purchaseInvoiceNo
		entry.DocumentNo = documentNo
		entry.RefundId = refundId
	} else {
		hashModel[key] = HashInvoiceEntry{
			Hash:       hash,
			NavID:      navId,
			InvoiceNo:  purchaseInvoiceNo,
			DocumentNo: documentNo,
			RefundId:   refundId,
		}
	}

	return nil
}

func UnmarshelCreateInvoiceResponse(stringData interface{}) (PostInvoiceEnvelope, error) {
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

func UnmarshelPostInvoiceResponse(stringData interface{}) (PostResponseInvoiceEnvelope, error) {
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

func InsertToNav(invoice WSPurchaseInvoicePage) (bool, error, interface{}) {
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Invoice.Sync.URL

	var result interface{}
	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(invoice)
	if err != nil {
		message := "Failed:InsertInvoice:Sync:3 Error mapping to XML -> " + err.Error()
		return false, errors.New(message), result
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
	isSuccess := false
	result, err = manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		message := "Failed:Sync:4 " + err.Error()
		return isSuccess, errors.New(message), result
	} else {
		resultStr, ok := result.(string)
		if !ok {
			// The type assertion failed
			message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
			return isSuccess, errors.New(message), result
		}
		match := utils.MatchRegexExpression(resultStr, `<WSPurchaseInvoicePage[^>]*>`)
		matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

		// Print the result
		if !match && matchFault {
			message := fmt.Sprintf("Failed:Sync:6 XML string does not contain <WSPurchaseInvoicePage> element: ", result)
			return isSuccess, errors.New(message), result
		} else {
			isSuccess = true
		}
	}
	return isSuccess, nil, result
}

func PostToNavAfterInsert(envelope PostInvoiceEnvelope) (bool, error, interface{}) {
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Invoice.Post.URL

	var result interface{}
	isSuccess := false
	// Map Go struct to XML

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
		// The type assertion failed
		message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
		return isSuccess, errors.New(message), result
	} else {
		utils.Console(result)
		resultPostStr, ok := result.(string)
		if !ok {
			// The type assertion failed
			message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
			return isSuccess, errors.New(message), result
		}
		match := utils.MatchRegexExpression(resultPostStr, `<PostPurchaseInvoice_Result[^>]*>`)
		matchFault := utils.MatchRegexExpression(resultPostStr, `<faultcode[^>]*>`)

		// Print the result
		if !match && matchFault {
			message := fmt.Sprintf("Failed:Sync:6 XML string does not contain <PostPurchaseInvoice_Result> element: ", result)
			return isSuccess, errors.New(message), result
		} else {
			isSuccess = true
		}
	}

	return isSuccess, nil, result
}

func SaveHashLogs(model HashInvoiceModel) (string, error) {
	//Paths
	PENDING_FILE_PATH := utils.INVOICE_HASH_FILE_PATH
	HASH_DB := utils.INVOICE_HASH_DB

	//Convert to String
	response, _ := data_parser.ParseModelToString(model)

	//Save to pending file
	var result string
	err := filesystem.CleanAndSave(PENDING_FILE_PATH, HASH_DB, response)
	if err != nil {
		message := "Failed:SaveHashLogs:Fetch:1 " + err.Error()
		//utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_FAILURE, timestamp+".json", message, "")
		return result, errors.New(message)
	} else {
		result := "Fetch: Successfully saved vendor to hashed file"
		//utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, PENDING_LOG_FILE_PATH, PENDING_SUCCESS, timestamp+".json", message, "")
		return result, nil
	}
}
