package providers

import "net/http"

type AtlassianClient struct {
	Client
}

func NewAtlassianClient(baseUrl, token string) *AtlassianClient {
	return &AtlassianClient{
		Client: Client{
			Client:  &http.Client{},
			BaseUrl: baseUrl,
			Token:   token,
		},
	}
}
