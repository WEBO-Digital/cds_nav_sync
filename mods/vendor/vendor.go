package vendor

import (
	"encoding/json"
	"fmt"
	"log"
	"nav_sync/config"
	filesystem "nav_sync/mods/afile_system"
	"nav_sync/mods/amanager"
	navapi "nav_sync/mods/anav_api"
	normalapi "nav_sync/mods/anormal_api"
	data_parser "nav_sync/mods/aparser"

	"nav_sync/utils"
)

func Fetch() {
	//Path
	VENDOR_FETCH_URL := config.Config.Vendor.Fetch.URL
	VENDOR_PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH

	//Fetch vendor data
	response, err := amanager.Fetch(VENDOR_FETCH_URL, normalapi.GET)
	if err != nil {
		utils.Console(err.Error())
	}
	utils.Console(response)

	//Save to pending file
	err = filesystem.Save(VENDOR_PENDING_FILE_PATH, response)
	if err != nil {
		utils.Console(err.Error())
	}
	utils.Console("Successfully saved vendor to pending file")
}

// func Sync() {
// 	//Eg.
// 	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
// 	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
// 	url := config.Config.Vendor.Sync.URL
// 	xmlPayload := `
// 		<Envelope
// 			xmlns="http://schemas.xmlsoap.org/soap/envelope/">
// 			<Body>
// 				<Create
// 					xmlns="urn:microsoft-dynamics-schemas/page/wsvendor">
// 					<WSVendor>
// 						<Name>Suman Intuji </Name>
// 						<Address>From vs code</Address>
// 						<Weighbridge_Supplier_ID>INJ123</Weighbridge_Supplier_ID>
// 					</WSVendor>
// 				</Create>
// 			</Body>
// 		</Envelope>
// 	`

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
	VENDOR_PENDING_FILE_PATH := utils.VENDOR_PENDING_FILE_PATH
	VENDOR_DONE_FILE_PATH := utils.VENDOR_DONE_FILE_PATH
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(VENDOR_PENDING_FILE_PATH)
	if err != nil {
		utils.Console(err.Error())
	}

	utils.Console(fileNames)

	if fileNames == nil || len(fileNames) < 1 {
		return
	}

	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//We assume here that Data are pushed to NAV

		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(VENDOR_PENDING_FILE_PATH, fileNames[i])

		// Step 2: Unmarshal JSON to struct
		var vendor WSVendor
		if err := json.Unmarshal([]byte(jsonData), &vendor); err != nil {
			utils.Console("Error unmarshaling JSON:", err)
		}

		//utils.Console(vendor)

		// Map Go struct to XML
		xmlData, err := data_parser.ParseJsonToXml(vendor)
		if err != nil {
			utils.Fatal("Error mapping to XML: ", err)
		}

		//Add XML envelope and body elements
		buffer := fmt.Sprintf(
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

		//Return the result
		//envolpeData := buffer.Bytes()
		xmlPayload := buffer //buffer.String()
		//utils.Console(xmlPayload)
		log.Println(xmlPayload)
		utils.Console("username: ", NTLM_USERNAME)
		utils.Console("username: ", NTLM_PASSWORD)
		utils.Console("URL: ", url)

		//Sync to Nav
		isSuccess := false
		result, err := amanager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
		if err != nil {
			utils.Console(err)
			isSuccess = false
		} else {
			utils.Console(result)
			isSuccess = true
		}

		if isSuccess {
			//Move to done file
			err = filesystem.MoveFile(fileNames[i], VENDOR_PENDING_FILE_PATH, VENDOR_DONE_FILE_PATH)
			if err != nil {
				utils.Console(err.Error())
			} else {
				utils.Console("File moved successfully")
			}
		}
	}
}
