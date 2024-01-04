package normalapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	data_parser "nav_sync/mods/ahelpers/parser"
	"net/http"
)

func Post(url string, token string, payloadData interface{}) (interface{}, error) {
	// Convert data to JSON
	jsonData, err := json.Marshal(payloadData)
	if err != nil {
		return nil, err
	}

	// // Make POST request
	// response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	return nil, err
	// }
	// defer response.Body.Close()

	// Create a new POST request with the JSON body
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Add the token to the request header
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json") // Set content type as JSON

	// Make the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		// Convert to JSON
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.New("unexpected status code: " + response.Status)
		}
		bodyString := string(body)
		return nil, errors.New(bodyString)
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Convert to JSON
	data, err := data_parser.ParseByteToJson(body)
	if err != nil {
		return nil, err
	}

	resData := data["payload"].(interface{})
	return resData, nil
}

func Get(url string, token string) (interface{}, error) {
	// //Make response call
	// Create a new request with the GET method
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add the token to the request header
	request.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	response, err := http.DefaultClient.Do(request)
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

	resData := data["payload"].(interface{})
	return resData, nil
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
