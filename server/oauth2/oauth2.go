package oauth2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	Author      string `json:"author,omitempty"`
}

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_id"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"state"`
	Author       string `json:"author"`
}

func AuthServer(w http.ResponseWriter, r *http.Request) {
	response := &TokenResponse{
		AccessToken: "theaccesstoken",
		Scope:       "thescope",
		TokenType:   "thetokentype",
	}

	resBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(resBytes))
}
