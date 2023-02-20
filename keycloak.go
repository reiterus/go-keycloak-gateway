package main

import (
	"encoding/json"
)

// TokenResponse structure
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// parseTokenResponse convert keycloak response to structure
func (su TokenResponse) parseTokenResponse(response string) TokenResponse {
	jsonData := []byte(response)
	result := TokenResponse{}
	_ = json.Unmarshal(jsonData, &result)

	return result
}
