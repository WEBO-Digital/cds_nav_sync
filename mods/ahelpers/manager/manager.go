package manager

import (
	"errors"
	"fmt"
	filesystem "nav_sync/mods/ahelpers/file_system"
	navapi "nav_sync/mods/ahelpers/nav_api"
	normalapi "nav_sync/mods/ahelpers/normal_api"
)

func Fetch(url string, method normalapi.APIMethod, tokenKey string, data interface{}) (interface{}, error) {
	//Make response call
	if method == normalapi.POST {
		return normalapi.Post(url, tokenKey, data)
	} else if method == normalapi.GET {
		return normalapi.Get(url, tokenKey)
	} else if method == normalapi.PUT {
		return normalapi.Put()
	} else {
		err := fmt.Sprintf("%s%s", "Invalid method found: ", method)
		return nil, errors.New(err)
	}
}

func Sync(url string, method navapi.NavMethod, xmlPayload string, user string, password string) (interface{}, error) {
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

func ApiFakeResponse(filePath string, fileName string) (bool, error, interface{}) {
	var result interface{}
	fakeByte, err := filesystem.ReadFile(filePath, fileName)
	if err != nil {
		return false, nil, result
	}
	result = string(fakeByte)
	return true, nil, result
}
