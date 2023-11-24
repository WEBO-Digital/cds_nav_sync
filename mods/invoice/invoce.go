package invoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"nav_sync/config"
	filesystem "nav_sync/mods/afile_system"
	"nav_sync/mods/amanager"
	"nav_sync/utils"
)

func Fetch() {
	//Path
	INVOICE_FETCH_URL := config.Config.Invoice.Fetch.URL
	INVOICE_PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH

	//Fetch vendor data
	response, err := amanager.Fetch(INVOICE_FETCH_URL)
	if err != nil {
		utils.Console(err.Error())
	}
	utils.Console(response)

	//Save to pending file
	err = filesystem.Save(INVOICE_PENDING_FILE_PATH, response)
	if err != nil {
		utils.Console(err.Error())
	}
	utils.Console("Successfully saved invoice to pending file")
}

func Sync() {
	//How to access pending directory
	//Case 1. Do this function run independently
	//Case 2. Do we run this function after fetch function is called

	//We go through Case 1.
	//Fetch all the pending files
	//Then sync one by one
	//After sync one file then move it to done folder.

	//Path
	INVOICE_PENDING_FILE_PATH := utils.INVOICE_PENDING_FILE_PATH
	INVOICE_DONE_FILE_PATH := utils.INVOICE_DONE_FILE_PATH

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(INVOICE_PENDING_FILE_PATH)
	if err != nil {
		utils.Console(err.Error())
	}

	//mods.Console(fileNames)

	if fileNames == nil || len(fileNames) < 1 {
		return
	}

	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//We assume here that Data are pushed to NAV

		//Get Json data from the file
		jsonData, err := filesystem.ReadFile(INVOICE_PENDING_FILE_PATH, fileNames[i])

		// Step 2: Unmarshal JSON to struct
		var invoice AddInvoiceModel
		if err := json.Unmarshal([]byte(jsonData), &invoice); err != nil {
			utils.Console("Error unmarshaling JSON:", err)
		}

		//utils.Console(invoice)

		// Map Go struct to XML
		xmlData, err := amanager.ParseJsonToXml(invoice)
		if err != nil {
			utils.Fatal("Error mapping to XML: ", err)
		}

		//Add XML envelope and body elements
		var buffer bytes.Buffer
		buffer.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
		buffer.WriteString(`<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">`)
		buffer.WriteString(`<Body>`)
		buffer.WriteString(`<Create xmlns="urn:microsoft-dynamics-schemas/page/wspurchaseinvoicepage">`)
		buffer.Write(xmlData)
		buffer.WriteString(`</Create>`)
		buffer.WriteString(`</Body>`)
		buffer.WriteString(`</Envelope>`)

		//Return the result
		envolpeData := buffer.Bytes()
		fmt.Println(string(envolpeData))

		//Sync to Nav

		//Move to done file
		err = filesystem.MoveFile(fileNames[i], INVOICE_PENDING_FILE_PATH, INVOICE_DONE_FILE_PATH)
		if err != nil {
			utils.Console(err.Error())
		}
	}

}
