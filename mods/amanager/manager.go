package amanager

import (
	"errors"
	"fmt"
	navapi "nav_sync/mods/anav_api"
	normalapi "nav_sync/mods/anormal_api"
	"nav_sync/utils"
)

func Fetch(url string, method normalapi.APIMethod) (interface{}, error) {
	// //Make response call
	// response, err := http.Get(url)
	// if err != nil {
	// 	return nil, err
	// }
	// defer response.Body.Close()

	// //Read the response body
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// //convert to json
	// data, err := ParseByteToJson(body)
	// if err != nil {
	// 	return nil, err
	// }

	if method == normalapi.POST {
		return normalapi.Post()
	} else if method == normalapi.GET {
		return normalapi.Get(url)
	} else if method == normalapi.PUT {
		return normalapi.Put()
	} else {
		err := fmt.Sprintf("%s%s", "Invalid method found: ", method)
		return nil, errors.New(err)
	}
}

func Sync(url string, method navapi.NavMethod, xmlPayload string, user string, password string) (interface{}, error) {
	utils.Console("method: ", method)
	if method == navapi.POST {
		return navapi.Post(url, xmlPayload, user, password)
	} else if method == navapi.GET {
		return navapi.Get()
	} else if method == navapi.PUT {
		return navapi.Put()
	} else {
		err := fmt.Sprintf("%s%s", "Invalid method found: ", method)
		return nil, errors.New(err)
	}
}
