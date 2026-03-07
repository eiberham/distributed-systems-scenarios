package clients

import (
	"encoding/json"
	"net/http"
)

type KeycloakClient struct {
	BaseURL string
	Realm   string
	Token   string
}

func NewKeycloakClient(baseURL, realm, token string) *KeycloakClient {
	return &KeycloakClient{
		BaseURL: baseURL,
		Realm:   realm,
		Token:   token,
	}
}

func (kc *KeycloakClient) GetUsers() (map[string]interface{}, error) {
	url := kc.BaseURL + "/admin/realms/" + kc.Realm + "/users"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+kc.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
