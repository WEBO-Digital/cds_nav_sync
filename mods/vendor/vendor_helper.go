package vendor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"nav_sync/config"
	filesystem "nav_sync/mods/ahelpers/file_system"
	"nav_sync/mods/ahelpers/manager"
	navapi "nav_sync/mods/ahelpers/nav_api"
	data_parser "nav_sync/mods/ahelpers/parser"

	"nav_sync/utils"
)

func GetHash() HashVendorModel {
	HASH_FILE_PATH := utils.VENDOR_HASH_FILE_PATH
	golbalModel := make(HashVendorModel)

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(HASH_FILE_PATH)
	if err != nil {
		//message := "Failed:Vendor:getHash:1 " + err.Error()
		return golbalModel //, errors.New(message)
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
			//message := "Failed:Vendor:getHash:2 " + err.Error()
			return golbalModel //, errors.New(message)
		}

		jsonString := string(jsonData)
		utils.Console("************************GetHash::inner****************************")
		utils.Console(jsonString)

		// Unmarshal JSON to struct
		localModel := make(HashVendorModel)
		if err := json.Unmarshal([]byte(jsonData), &localModel); err != nil {
			//message := "Failed:Sync:2 Error unmarshaling JSON -> " + err.Error()
			return golbalModel //, errors.New(message)
		}

		// Append localModel to golbalModel
		for key, value := range localModel {
			golbalModel[key] = value
		}
	}

	return golbalModel //, nil
}

func CompareWithHash(vendor WSVendor, hash HashVendorModel) (HashVendorModel, bool) {
	//var hashRecord HashVendorModel
	hashRecord := hash
	vendorStr, err := data_parser.ParseModelToString(vendor)
	if err != nil {
		//message := "Failed:Vendor:getHash:2 " + err.Error()
		//return golbalModel //, errors.New(message)
		return hashRecord, true
	}
	hashStr := utils.ComputeMD5(vendorStr)
	for key, _ := range hash {
		//hash is already pushed to nav server
		if containsHash(hash, hashStr) {
			//hashRecord[key] = value
			return hashRecord, true
		}

		if !containsHash(hash, hashStr) && key == vendor.WeighbridgeSupplierID { //&& value.NavID == nil
			/** @TODO: UPDATE THE VEDOR **/
			return hashRecord, true
		}
	}

	//hash is not found then insert to nav server
	return hashRecord, false
}

func containsHash(model HashVendorModel, targetHash string) bool {
	for _, entry := range model {
		if entry.Hash == targetHash {
			return true
		}
	}
	return false
}

func UpdateHashInModel(hashModel HashVendorModel, key, hash string, navId string) error {
	// Initialize the map if it is nil
	if hashModel == nil {
		hashModel = make(HashVendorModel)
	}

	// Check if the key exists in the map
	if entry, exists := hashModel[key]; exists {
		// Update the Hash field with the new value
		entry.Hash = hash
		entry.NavID = &navId
	} else {
		hashModel[key] = HashVendorEntry{
			Hash:  hash,
			NavID: &navId,
		}
	}

	return nil
}

func InsertToNav(vendor WSVendor, fileName string, jsonString string) (bool, error, interface{}) {
	//Path
	NTLM_USERNAME := config.Config.Auth.Ntlm.Username
	NTLM_PASSWORD := config.Config.Auth.Ntlm.Password
	url := config.Config.Vendor.Sync.URL

	var result interface{}
	// Map Go struct to XML
	xmlData, err := data_parser.ParseJsonToXml(vendor)
	if err != nil {
		message := "Failed:Sync:3 Error mapping to XML -> " + err.Error()
		//utils.Console(message)
		//logger.LogNavState(logger.SUCCESS, DONE_LOG_FILE_PATH, DONE_FAILURE, fileName, message, jsonString)
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

	//Return the result
	log.Println(xmlPayload)
	utils.Console("username: ", NTLM_USERNAME)
	utils.Console("username: ", NTLM_PASSWORD)
	utils.Console("URL: ", url)

	//Sync to Nav
	isSuccess := false
	result, err = manager.Sync(url, navapi.POST, xmlPayload, NTLM_USERNAME, NTLM_PASSWORD)
	if err != nil {
		message := "Failed:Sync:4 " + err.Error()
		// utils.Console(message)
		// logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileName, message, xmlPayload)
		return false, errors.New(message), result
	} else {
		resultStr, ok := result.(string)
		if !ok {
			// The type assertion failed
			message := fmt.Sprintf("Failed:Sync:5 Could not convert to string: ", result)
			// utils.Console(message)
			// logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileName, message, xmlPayload)
			return false, errors.New(message), result
		}
		match := utils.MatchRegexExpression(resultStr, `<Create_Result[^>]*>`)
		matchFault := utils.MatchRegexExpression(resultStr, `<faultcode[^>]*>`)

		// Print the result
		if !match && matchFault {
			message := fmt.Sprintf("Failed:Sync:6 XML string does not contain <Create_Result> element: ", result)
			// utils.Console(message)
			// logger.LogNavState(logger.FAILURE, DONE_LOG_FILE_PATH, DONE_FAILURE, fileName, message, resultStr)
			return false, errors.New(message), result
		} else {
			isSuccess = true
		}
	}

	return isSuccess, nil, result
}

func SaveHashLogs(model HashVendorModel) (string, error) {
	//Paths
	PENDING_FILE_PATH := utils.VENDOR_HASH_FILE_PATH
	HASH_DB := utils.VENDOR_HASH_DB

	//Convert to String
	response, _ := data_parser.ParseModelToString(model)

	//Save to pending file
	var result string
	err := filesystem.CleanAndSave(PENDING_FILE_PATH, HASH_DB, response)
	if err != nil {
		message := "Failed:Fetch:2 " + err.Error()
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
