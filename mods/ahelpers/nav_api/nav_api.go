package navapi

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/azure/go-ntlmssp"
)

func Post(url string, xmlPayload string, user string, password string) (interface{}, error) {
	//Make response call
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	//Create the payload
	payload := bytes.NewBuffer([]byte(xmlPayload))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Set("Soapaction", "urn:microsoft-dynamics-schemas/page/wsvendor:Create")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	//Do call
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//Get response
	return string(body), nil
}

func Get() (interface{}, error) {
	return nil, nil
}

func Put() (interface{}, error) {
	return nil, nil
}

func Delete() (interface{}, error) {
	return nil, nil
}

// Creating different types of methods
type NavMethod string

const (
	POST   NavMethod = "POST"
	GET    NavMethod = "GET"
	PUT    NavMethod = "PUT"
	DELETE NavMethod = "DELETE"
)
