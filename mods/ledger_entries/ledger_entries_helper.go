package ledgerentries

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

func UnmarshalStringToLedgerEntries(stringData interface{}) ([]LedgerEntriesCreate, error) {
	var ledgers []LedgerEntriesCreate
	// Type assertion to get the underlying string
	str, ok := stringData.(string)
	if !ok {
		return ledgers, errors.New("UnmarshalStringToLedgerEntries: Conversion failed: not a string")
	}

	// Convert the string to a byte slice
	xmlData := []byte(str)

	// Map Go struct to XML
	err := xml.Unmarshal(xmlData, &ledgers)
	if err != nil {
		return ledgers, errors.New("UnmarshalStringToLedgerEntries: Error decoding XML: " + err.Error())
	}
	return ledgers, nil
}

func UnmarshelCreateLedgerEntryResponse(stringData interface{}) (PostLedgerEntriesEnvelope, error) {
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

func InsertToNav(ledgerentrie LedgerEntriesCreate) (bool, error, interface{}) {
	//Path
	FAKE_INSERT := config.Config.LedgerEntries.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.LedgerEntries.Sync.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Post To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "ledger_entries_insert_fake.xml")
		return isFakeSuccess, err, result
	}

	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(ledgerentrie.VendorPayment)
	if err != nil {
		return false, errors.New("insertLedgerEntries: Error mapping to XML -> " + err.Error()), result
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
		ledgerentrie.CurrentJnlBatchName,
		string(xmlData),
	)

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
		match := utils.MatchRegexExpression(resultStr, `<VendorPayment[^>]*>`)
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

func PostLedgerEntriesAfterCreation(envelope PostLedgerEntriesEnvelope) (bool, error, interface{}) {
	//Path
	FAKE_INSERT := config.Config.LedgerEntries.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.LedgerEntries.Post.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Post To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "ledger_entries_post_fake.xml")
		return isFakeSuccess, err, result
	}

	isSuccess := false
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

	//Sync to Nav
	result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		// The type assertion failed
		message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
		return isSuccess, errors.New(message), result
	} else {
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
