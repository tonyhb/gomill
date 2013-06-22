package gomill

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
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
	post() (response interface{}, urlResource string)
}

// Create a new Paymill API struct with your API Key
func New(apiKey string) PaymillAPI {
	return PaymillAPI{
		Key: apiKey,
	}
}

func (this *PaymillAPI) Create(r Resource) (response interface{}, err error) {
	var urlResource string
	var urlValues url.Values
	// Find out which type of request we're making. We've been passed a
	// pointer to a resource and haven't used reflect.Indirect to get the
	// value, so the names are prefixed with a * below.
	switch reflect.TypeOf(r).String() {
	case "*gomill.Transaction":
		response, urlResource = r.post()
		urlValues = url.Values{
			"amount":      {r.(*Transaction).Amount},
			"currency":    {r.(*Transaction).Currency},
			"token":       {r.(*Transaction).Token},
			"description": {r.(*Transaction).Description.(string)},
		}
	}
	// Make the request to paymill
	resp, err := http.PostForm("https://"+this.Key+":@api.paymill.com/v2/"+urlResource, urlValues)
	if err != nil {
		return new(interface{}), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return new(interface{}), err
	}
	// Check to see if this is an error.
	if resp.StatusCode >= 400 {
		response = new(PaymillError)
		err = errors.New("Paymill API Error: Status " + resp.Status)
	}
	json.Unmarshal(body, response)
	return
}

