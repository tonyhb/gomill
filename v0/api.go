package gomill

const (
	API_VERSION = "v2"
)

type PaymillAPI struct {
	Key string
}

// Create a new Paymill API struct with your API Key
func New(apiKey string) PaymillAPI {
	return PaymillAPI{
		Key: apiKey,
	}
}
