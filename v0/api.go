package gomill

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PaymillAPI struct {
	Key string
}

// Paymill returns a JSON error along with a status code
type PaymillError struct {
	Error     interface{} `json:"error"`     // This is sometimes a string and sometimes a map[string]interface{}
	Exception interface{} `json:"exception"` // This is sometimes  astring and sometimes not provided
}

type Resource interface {
	create() (response interface{}, urlResource string, urlParams url.Values)
}

// Create a new Paymill API struct with your API Key
func New(apiKey string) PaymillAPI {
	return PaymillAPI{
		Key: apiKey,
	}
}

func (this *PaymillAPI) Create(r Resource) (response interface{}, err error) {
	// Make the request to paymill
	response, urlResource, urlParams := r.create()
	resp, err := http.PostForm("https://"+this.Key+":@api.paymill.com/v2/"+urlResource, urlParams)
	if err != nil {
		return new(interface{}), err
	}
	// Read the contents of the body into a bytes buffer
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return new(interface{}), err
	}
	// Check to see if this is an error so we can unmarshal into a PaymillError struct
	if resp.StatusCode >= 400 {
		response = new(PaymillError)
		err = errors.New("Paymill API Error: Status " + resp.Status)
	}
	// Unmarshal into either the struct provided by (r Resource) or a PaymillError
	json.Unmarshal(body, response)
	return
}

