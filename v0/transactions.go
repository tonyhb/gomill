package gomill

import (
	"reflect"
	"strconv"
)

type Transaction struct {
	Id               string        `json:"id"`
	Amount           string        `json:"amount"`
	Currency         string        `json:"currency"`
	Token            string        `json:"token"` // Set when creating a transaction, "" otherwise
	Description      interface{}   `json:"description"`
	OriginAmount     int           `json:"origin_amount"`
	Status           string        `json:"status"`
	Livemode         bool          `json:"livemode"`
	Refunds          interface{}   `json:"refunds"`
	CreatedAt        int           `json:"created_at"`
	UpdatedAt        int           `json:"updated_at"`
	ResponseCode     int           `json:"response_code"`
	ShortId          interface{}   `json:"short_id"`
	Invoices         []interface{} `json:"invoices"`
	Fees             []interface{} `json:"fees"`
	Payment          Payment       `json:"payment"`
	Client           Client        `json:"client"`
	Preauthorization interface{}   `json:"preauthorization"`
	Mode             interface{}   `json:"mode"`
}

// Sets the Amount field to a given value and optionally multiplies by 100 if
// the unit isn't minor (eg. dollars instead of cents)
// Amount can be in the int, uint, float or string family.
func (this *Transaction) SetAmount(amount interface{}, inCents bool) {
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
	this.Amount = value
}

func (this *Transaction) post() (response interface{}, urlResource string) {
	return new(Transaction), "transactions"
}
