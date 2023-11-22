package amanager

import "encoding/json"

func ParseToJson(body []byte) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
