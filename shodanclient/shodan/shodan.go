package shodan

const BaseURL = "http://api.shodan.io"

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
