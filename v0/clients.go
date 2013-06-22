package gomill

type Client struct {
	Id           string      `json:"id"`
	Email        interface{} `json:"email"`
	Description  interface{} `json:"description"`
	CreatedAt    int         `json:"created_at"`
	UpdatedAt    int         `json:"updated_at"`
	Payment      []Payment   `json:"payment"`
	Subscription interface{} `json:"subscription"` // @TODO add subscriptions struct
}
