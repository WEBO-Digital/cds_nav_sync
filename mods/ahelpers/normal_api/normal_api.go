package normalapi

import (
	"io/ioutil"
	data_parser "nav_sync/mods/ahelpers/parser"
	"net/http"
)

func Post() (interface{}, error) {
	return nil, nil
}

func Get(url string) (interface{}, error) {
	//Make response call
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

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
