package clients

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type KeycloakClient struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

func NewKeycloakClient(baseURL string, realm string, clientID string, clientSecret string) *KeycloakClient {
	return &KeycloakClient{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (kc *KeycloakClient) GetUsers() ([]map[string]interface{}, error) {
	url := kc.BaseURL + "/admin/realms/" + kc.Realm + "/users"
	req, _ := http.NewRequest("GET", url, nil)
	token, err := kc.GetServiceToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (k *KeycloakClient) GetServiceToken() (string, error) {
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {k.ClientID}, // users-service
		"client_secret": {k.ClientSecret},
	}

	resp, _ := http.PostForm(
		k.BaseURL+"/realms/"+k.Realm+"/protocol/openid-connect/token",
		data,
	)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result["access_token"].(string), nil
}
