package vendor

import (
	"encoding/json"
	"fmt"
	"log"
	"nav_sync/config"
	"nav_sync/logger"
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
	VENDOR_PENDING_LOG_FILE_PATH := utils.VENDOR_DONE_LOG_FILE_PATH
	VENDOR_PENDING_FAILURE := utils.VENDOR_PENDING_FAILURE
	VENDOR_PENDING_SUCCESS := utils.VENDOR_PENDING_SUCCESS

	//Fetch vendor data
	response, err := amanager.Fetch(VENDOR_FETCH_URL, normalapi.GET)
	if err != nil {
		message := "Failed:Fetch:1 " + err.Error()
		utils.Console(message)
		logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_PENDING_LOG_FILE_PATH, VENDOR_PENDING_FAILURE, "", message, "")
	}
	utils.Console(response)

	//get current timestamp
	timestamp := utils.GetCurrentTime()

	//Save to pending file
	err = filesystem.Save(VENDOR_PENDING_FILE_PATH, timestamp+".json", response)
	if err != nil {
		message := "Failed:Fetch:2 " + err.Error()
		utils.Console(message)
		logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_PENDING_LOG_FILE_PATH, VENDOR_PENDING_FAILURE, timestamp+".json", message, "")
	} else {
		message := "Fetch: Successfully saved vendor to pending file"
		utils.Console(message)
		logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_PENDING_LOG_FILE_PATH, VENDOR_PENDING_SUCCESS, timestamp+".json", message, "")
	}

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
	VENDOR_DONE_LOG_FILE_PATH := utils.INVOICE_DONE_LOG_FILE_PATH
	VENDOR_DONE_FAILURE := utils.INVOICE_DONE_FAILURE
	VENDOR_DONE_SUCCESS := utils.INVOICE_DONE_SUCCESS
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(VENDOR_PENDING_FILE_PATH)
	if err != nil {
		message := "Failed:Sync:1 " + err.Error()
		utils.Console(message)
		logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_DONE_LOG_FILE_PATH, VENDOR_DONE_FAILURE, "", message, "")
	}

	utils.Console(fileNames)

	if fileNames == nil || len(fileNames) < 1 {
		return
	}

	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(VENDOR_PENDING_FILE_PATH, fileNames[i])

		// Unmarshal JSON to struct
		var vendor WSVendor
		if err := json.Unmarshal([]byte(jsonData), &vendor); err != nil {
			message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
			utils.Console(message)
			logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_DONE_LOG_FILE_PATH, VENDOR_DONE_FAILURE, fileNames[i], message, "")
		}

		//utils.Console(vendor)

		// Map Go struct to XML
		xmlData, err := data_parser.ParseJsonToXml(vendor)
		if err != nil {
			message := "Failed:Sync:3 Error mapping to XML -> " + err.Error()
			utils.Console(message)
			logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_DONE_LOG_FILE_PATH, VENDOR_DONE_FAILURE, fileNames[i], message, "")
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

		//Return the result
		log.Println(xmlPayload)
		utils.Console("username: ", NTLM_USERNAME)
		utils.Console("username: ", NTLM_PASSWORD)
		utils.Console("URL: ", url)

		//Sync to Nav
		isSuccess := false
		result, err := amanager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
		if err != nil {
			isSuccess = false
			message := "Failed:Sync:4 " + err.Error()
			utils.Console(message)
			logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_DONE_LOG_FILE_PATH, VENDOR_DONE_FAILURE, fileNames[i], message, "")
		} else {
			utils.Console(result)
			isSuccess = true
		}

		if isSuccess {
			//Move to done file
			err = filesystem.MoveFile(fileNames[i], VENDOR_PENDING_FILE_PATH, VENDOR_DONE_FILE_PATH)
			if err != nil {
				message := "Failed:Sync:5 " + err.Error()
				utils.Console(message)
				logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_DONE_LOG_FILE_PATH, VENDOR_DONE_FAILURE, fileNames[i], message, "")
			} else {
				message := "File moved successfully"
				utils.Console(message)
				logger.LogInvoiceFetch(logger.SUCCESS, VENDOR_DONE_LOG_FILE_PATH, VENDOR_DONE_SUCCESS, fileNames[i], message, "")
			}
		}
	}
}
