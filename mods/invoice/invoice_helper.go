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

func UnmarshalStringToInvoice(stringData interface{}) ([]WSPurchaseInvoicePage, error) {
	var invoices []WSPurchaseInvoicePage

	jsonData, err := json.Marshal(stringData)
	if err != nil {
		return invoices, errors.New("Conversion failed: " + err.Error())
	}

	// Map Go struct to XML
	err = json.Unmarshal(jsonData, &invoices)
	if err != nil {
		return invoices, errors.New("UnmarshalStringToInvoice: Error decoding XML: " + err.Error())
	}
	return invoices, nil
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
	//Path
	FAKE_INSERT := config.Config.Invoice.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Invoice.Sync.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Insert To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "invoice_insert_fake.xml")
		return isFakeSuccess, err, result
	}

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

	err = filesystem.Save("/data/invoice/test/", utils.GetCurrentTime(), xmlPayload)
	if err != nil {
		utils.Console(err.Error())
	}

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
	//Path
	FAKE_INSERT := config.Config.Invoice.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Invoice.Post.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Post To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "invoice_post_fake.xml")
		return isFakeSuccess, err, result
	}

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

	//Sync to Nav
	result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		// The type assertion failed
		message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
		return isSuccess, errors.New(message), result
	} else {
		//utils.Console(result)
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
