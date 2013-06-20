package gomill

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

type TransactionRequest struct {
	Id          int    `paymill:"-"`
	Amount      string `paymill:"amount"`
	Currency    string `paymill:"currency"`
	Token       string `paymill:"token"`
	Description string `paymill:"description,omitempty"`
}

type TransactionResponse struct {
	Data struct {
		Id            string
		Amount        string
		Origin_amount int
		Status        string
		Description   interface{}
		Livemode      bool
		Refunds       interface{}
		Currency      string
		Created_at    int
		Updated_at    int
		Response_code int
		Invoices      []struct {
		}
		Payment struct {
			Id           string
			Type         string
			Client       string
			Card_type    string
			Country      interface{}
			Expire_month int
			Expire_year  int
			Card_holder  interface{}
			Lsat4        string
			Created_at   int
			Updated_at   int
		}
		Client struct {
			Id           string
			Email        interface{}
			Description  interface{}
			Created_at   int
			Updated_at   int
			Parment      []struct{}
			Subscription interface{}
		}
		Preauthorization interface{}
	}
	Error     interface{}
	Exception interface{}
	Mode      interface{}
}

// Sets the Amount field to a given value and optionally multiplies by 100 if
// the unit isn't minor (eg. dollars instead of cents)
// Amount can be in the int, uint, float or string family.
func (req *TransactionRequest) SetAmount(amount interface{}, inCents bool) {
	var value string
	switch amount.(type) {
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(amount).Int()
		if inCents == false {
			v = v * 100
		}
		value = strconv.FormatInt(v, 10)
	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(amount).Uint()
		if inCents == true {
			v = v * 100
		}
		value = strconv.FormatUint(v, 10)
	case float32, float64:
		v := reflect.ValueOf(amount).Float()
		if inCents == false {
			v = v * 100
		}
		value = strconv.FormatFloat(v, 'f', 0, 64)
	case string:
		// We need to convert to a float because the string may be something
		// like 100.23, meaning adding "00" to the end won't work.
		if inCents == false {
			v, _ := strconv.ParseFloat(reflect.ValueOf(amount).String(), 64)
			v = v * 100
			value = strconv.FormatFloat(v, 'f', 0, 64)
		} else {
			value = reflect.ValueOf(amount).String()
		}
	}
	req.Amount = value
}

func (req *TransactionRequest) Create(api PaymillAPI) (response *TransactionResponse, err error) {
	response = new(TransactionResponse)
	// Convert the request struct to a url values map and post to Paymill
	values := structToMap(req)
	resp, err := http.PostForm("https://"+api.Key+":@api.paymill.com/v2/transactions", values)
	if err != nil {
		return
	}
	// Take the paymill response, unmarshall it into a struct and check to see if we were successful in capturing the payment
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(body, response)
	return
}
