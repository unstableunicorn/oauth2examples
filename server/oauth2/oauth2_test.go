package oauth2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccessCode(t *testing.T) {
	t.Run("returns access code", func(t *testing.T) {
		reqBody := TokenRequest{
			ClientID:     "theclientid",
			ClientSecret: "thecliendsecret",
			Code:         "thecode",
			RedirectURI:  "theredirecturi",
			State:        "thestate",
			Author:       "github",
		}

		body, _ := json.Marshal(reqBody)

		request, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		AuthServer(response, request)

		want := TokenResponse{
			AccessToken: "theaccesstoken",
			Scope:       "thescope",
			TokenType:   "thetokentype",
			Author:      "github",
		}

		var got TokenResponse
		json.NewDecoder(response.Body).Decode(&got)

		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}

	})
}
