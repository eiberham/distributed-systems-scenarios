package providers

import "net/http"

type Client struct {
	Client  *http.Client
	BaseUrl string
	Token   string
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	return c.Client.Do(req)
}
