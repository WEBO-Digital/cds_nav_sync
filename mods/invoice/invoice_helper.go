package invoice

import (
	"encoding/xml"
	"errors"
	"fmt"
	"nav_sync/config"
	"nav_sync/mods/ahelpers/manager"
	navapi "nav_sync/mods/ahelpers/nav_api"
	data_parser "nav_sync/mods/ahelpers/parser"
	"nav_sync/utils"
)

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
