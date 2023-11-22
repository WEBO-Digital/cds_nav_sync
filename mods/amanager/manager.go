package amanager

import (
	"io/ioutil"
	"net/http"
)

func Fetch(url string) (interface{}, error) {
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
	data, err := ParseToJson(body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Sync() (interface{}, error) {
	return nil, nil
}
