package normalapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	data_parser "nav_sync/mods/ahelpers/parser"
	"net/http"
)

func Post(url string, data interface{}) (interface{}, error) {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Make POST request
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + response.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Convert to JSON
	var responseData interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func Get(url string) (interface{}, error) {
	//Make response call
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + response.Status)
	}

	//Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//convert to json
	data, err := data_parser.ParseByteToJson(body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Put() (interface{}, error) {
	return nil, nil
}

func Delete() (interface{}, error) {
	return nil, nil
}

// Creating different types of methods
type APIMethod string

const (
	POST   APIMethod = "POST"
	GET    APIMethod = "GET"
	PUT    APIMethod = "PUT"
	DELETE APIMethod = "DELETE"
)
