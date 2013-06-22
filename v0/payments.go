package gomill

type Payment struct {
	Id          string      `json:"id"`
	Type        string      `json:"type"`
	Code        string      `json:"code"`    // Direct debit only
	Account     string      `json:"account"` // Direct debit only
	Holder      string      `json:"holder"`  // Direct debit only
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
