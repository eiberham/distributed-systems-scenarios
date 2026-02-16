package providers

import "net/http"

type GitHubClient struct {
	Client
}

func NewGitHubClient(baseUrl, token string) *GitHubClient {
	return &GitHubClient{
		Client: Client{
			Client:  &http.Client{},
			BaseUrl: baseUrl,
			Token:   token,
		},
	}
}
