package vendor

import (
	"encoding/json"
	"errors"
	"fmt"
	"nav_sync/config"
	"nav_sync/mods/ahelpers/manager"
	navapi "nav_sync/mods/ahelpers/nav_api"
	data_parser "nav_sync/mods/ahelpers/parser"

	"nav_sync/utils"
)

func InsertToNav(vendor WSVendor) (bool, error, interface{}) {
	//Path
	FAKE_INSERT := config.Config.Vendor.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Insert To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "vendor_fake.xml")
		return isFakeSuccess, err, result
	}

	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(vendor)
	if err != nil {
		message := "Failed:InsertToNav:Sync:3 Error mapping to XML -> " + err.Error()
		return false, errors.New(message), result
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<Create xmlns="urn:microsoft-dynamics-schemas/page/wsvendor">
						%s
					</Create>
				</Body>
			</Envelope>
		`,
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
		match := utils.MatchRegexExpression(resultStr, `<Create_Result[^>]*>`)
		matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

		// Print the result
		if !match && matchFault {
			message := fmt.Sprintf("Failed:Sync:6 XML string does not contain <Create_Result> element: ", result)
			return isSuccess, errors.New(message), result
		} else {
			isSuccess = true
		}
	}

	return isSuccess, nil, result
}

func UpdateToNav(vendor WSVendor) (bool, error, interface{}) {
	//Path
	FAKE_INSERT := config.Config.Vendor.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Insert To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "vendor_get_key_fake.xml")
		return isFakeSuccess, err, result
	}

	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(vendor)
	if err != nil {
		message := "Failed:Update:Sync:3 Error mapping to XML -> " + err.Error()
		return false, errors.New(message), result
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<Update xmlns="urn:microsoft-dynamics-schemas/page/wsvendor">
						%s
					</Update>
				</Body>
			</Envelope>
		`,
		string(xmlData),
	)

	//Sync to Nav
	isSuccess := false
	result, err = manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		message := "Failed:Update:4 " + err.Error()
		return isSuccess, errors.New(message), result
	} else {
		resultStr, ok := result.(string)
		if !ok {
			// The type assertion failed
			message := fmt.Sprintf("Failed:Update:5 Could not convert to string: ", result)
			return isSuccess, errors.New(message), result
		}
		match := utils.MatchRegexExpression(resultStr, `<Update_Result[^>]*>`)
		matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

		// Print the result
		if !match && matchFault {
			message := fmt.Sprintf("Failed:Update:6 XML string does not contain <Update_Result> element: ", result)
			return isSuccess, errors.New(message), result
		} else {
			isSuccess = true
		}
	}

	return isSuccess, nil, result
}

func GetKeyFromVendorId(navId string, weighbridgeSupplierID string) (bool, error, interface{}) {
	//Path
	FAKE_INSERT := config.Config.Vendor.FakeInsert
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	//Result
	var result interface{}

	if FAKE_INSERT {
		//Fake Insert To Nav
		isFakeSuccess, err, result := manager.ApiFakeResponse("/ztest/", "vendor_get_key_fake.xml")
		return isFakeSuccess, err, result
	}

	//Add XML envelope and body elements
	xmlPayload := fmt.Sprintf(
		`
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<Read xmlns="urn:microsoft-dynamics-schemas/page/wsvendor">
						<No>%s</No>
						<Weighbridge_Supplier_ID>%s</Weighbridge_Supplier_ID>
					</Read>
				</Body>
			</Envelope>
		`,
		navId,
		weighbridgeSupplierID,
	)

	//Sync to Nav
	isSuccess := false
	result, err := manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		message := "Failed:GetKeyFromVendorId:4 " + err.Error()
		return isSuccess, errors.New(message), result
	} else {
		resultStr, ok := result.(string)
		if !ok {
			// The type assertion failed
			message := fmt.Sprintf("Failed:GetKeyFromVendorId:5 Could not convert to string: ", result)
			return isSuccess, errors.New(message), result
		}
		match := utils.MatchRegexExpression(resultStr, `<Read_Result[^>]*>`)
		//matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

		// Print the result
		// if !match && matchFault {
		if !match {
			message := fmt.Sprintf("Failed:GetKeyFromVendorId:6 XML string does not contain <Read_Result> element: ", result)
			return isSuccess, errors.New(message), result
		} else {
			isSuccess = true
		}
	}

	return isSuccess, nil, result
}

func UnmarshelByteToVendor(jsonData []byte) ([]WSVendor, error) {
	// Unmarshal JSON to struct
	var vendors []WSVendor

	if err := json.Unmarshal([]byte(jsonData), &vendors); err != nil {
		return vendors, err
	}

	return vendors, nil
}

func UnmarshalStringToVendor(stringData interface{}) ([]WSVendor, error) {
	var vendors []WSVendor
	jsonData, err := json.Marshal(stringData)
	if err != nil {
		return vendors, errors.New("Conversion failed: " + err.Error())
	}

	// Map Go struct to XML
	err = json.Unmarshal(jsonData, &vendors)
	if err != nil {
		return vendors, errors.New("unmarshelCreateInvoiceResponse: Error decoding XML: " + err.Error())
	}
	return vendors, nil
}
