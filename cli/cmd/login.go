/*Package cmd login module
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Will launch a webpage and provide login details",
	Long:  `Some longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		startLogin()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	ctx          context.Context
	clientId     = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	authUrl      = "http://localhost:3000"
	tokenUrl     = "https://github.com/login/oauth/access_token"
	ghUserApi    = "https://api.github.com/user"
	// my own callback URL
	redirectUrl = "http://localhost:3001/oauth/callback"
)

type gitHubAccessResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type gitHubUserDetails struct {
	User      string `json:"name"`
	Alias     string `json:"login"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	queryParts, _ := url.ParseQuery(r.URL.RawQuery)

	// Use the authorization code that is pushed to the redirect

	code := queryParts["code"][0]
	log.Printf("code: %s\n", code)

	// this needs to be performed by another server
	requestBodyMap := map[string]string{"client_id": clientId, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to github auth
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var ghar gitHubAccessResponse
	json.Unmarshal(body, &ghar)
	log.Println(ghar)

	// this in it's own thing as well
	user := getGitHubUserDetails(ghar)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", user)
}

func getGitHubUserDetails(ghar gitHubAccessResponse) gitHubUserDetails {
	// POST request to github auth
	req, err := http.NewRequest("GET", ghUserApi, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", ghar.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Got the body, generic unmarshal:")
	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal("Could not parse reponse body")
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("Failed to prettify data")
	}

	fmt.Println(string(b))

	fmt.Printf("Name: %s", data["name"].(string))
	fmt.Printf("Email: %s", data["email"].(string))

	fmt.Println("Got the body, specific unmarshal")
	var userDetails gitHubUserDetails

	if err := json.Unmarshal(body, &userDetails); err != nil {
		log.Fatal("Could not parse reponse body")
	}
	return userDetails
}

func startLogin() {
	ctx = context.Background()
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	state := uuid.New().String()
	localUrl, err := url.Parse(authUrl)
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	params.Add("state", state)
	params.Add("scope", "read:user read:email")
	params.Add("redirect_uri", redirectUrl)

	localUrl.RawQuery = params.Encode()
	log.Println(color.CyanString("You will now be taken to your browser for authentication"))
	time.Sleep(1 * time.Second)

	err = open.Start(localUrl.String())

	if err != nil {
		fmt.Println(color.BlueString("could auto launch webpage, please copy and past the following url:"))
		fmt.Println(localUrl.String())
	}

	time.Sleep(1 * time.Second)
	mux := http.NewServeMux()
	server := http.Server{Addr: ":3001", Handler: mux}
	mux.HandleFunc("/oauth/callback", callbackHandler)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		server.Shutdown(ctx)
		return
	}
}
