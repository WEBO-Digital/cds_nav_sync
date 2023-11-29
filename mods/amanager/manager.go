package amanager

import (
	"errors"
	"fmt"
	navapi "nav_sync/mods/anav_api"
	normalapi "nav_sync/mods/anormal_api"
	"nav_sync/utils"
)

func Fetch(url string, method normalapi.APIMethod) (interface{}, error) {
	//Make response call
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
	//Make response call
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
