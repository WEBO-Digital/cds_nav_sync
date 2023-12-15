package vendor

import (
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
