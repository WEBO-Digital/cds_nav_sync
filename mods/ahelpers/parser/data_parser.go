package data_parser

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"nav_sync/utils"
)

func ParseByteToJson(body []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseModelToString(body interface{}) (string, error) {
	var conv string
	// Convert struct to JSON
	jsonData, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return conv, err
	}

	// Print the JSON representation
	conv = string(jsonData)
	return conv, nil
}

func ParseJsonToXml(data interface{}) ([]byte, error) {
	//Convert Uint8 To String
	//data = convertUint8ToString(data)

	//Format into xml
	xmlData, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		utils.Console(err.Error())
		return nil, err
	}

	return xmlData, nil
}

func convertUint8ToString(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = convertUint8ToString(value)
		}
	case []interface{}:
		for i, item := range v {
			v[i] = convertUint8ToString(item)
		}
	case []uint8:
		return string(v)
	}

	return data
}
