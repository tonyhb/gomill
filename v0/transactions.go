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
		Id           string        `json:"id"`
		Amount       string        `json:"amount"`
		OriginAmount int           `json:"origin_amount"`
		Status       string        `json:"status"`
		Description  interface{}   `json:"description"`
		Livemode     bool          `json:"livemode"`
		Refunds      interface{}   `json:"refunds"`
		Currency     string        `json:"currency"`
		CreatedAt    int           `json:"created_at"`
		UpdatedAt    int           `json:"updated_at"`
		ResponseCode int           `json:"response_code"`
		ShortId      interface{}   `json:"short_id"`
		Invoices     []interface{} `json:"invoices"`
		Fees         []interface{} `json:"fees"`
		Payment      struct {
			Id          string      `json:"id"`
			Type        string      `json:"type"`
			Client      string      `json:"client"`
			CardType    string      `json:"card_type"`
			Country     interface{} `json:"country"`
			ExpireMonth interface{} `json:"expire_month"`
			ExpireYear  interface{} `json:"expire_year"`
			CardHolder  interface{} `json:"card_holder"`
			Last4       string      `json:"last4"`
			CreatedAt   int         `json:"created_at"`
			UpdatedAt   int         `json:"updated_at"`
		}
		Client struct {
			Id           string        `json:"id"`
			Email        interface{}   `json:"email"`
			Description  interface{}   `json:"description"`
			CreatedAt    int           `json:"created_at"`
			UpdatedAt    int           `json:"updated_at"`
			Payment      []interface{} `json:"payment"`
			Subscription interface{}   `json:"subscription"`
		}
		Preauthorization interface{} `json:"preauthorization"`
	}
	Error     interface{} `json:"error"`
	Exception interface{} `json:"exception"`
	Mode      interface{} `json:"mode"`
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
