package gocord

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Implements Oauth2 helper methods and definitions. This package cannot be used standalone, and requires a website or
// another app implementing Oauth flow directly.

// Oauth2Application represents the application interacting with the Oauth2 flow
type Oauth2Application struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

// Oauth2Callback is a struct of data returned in the querystring params during Oauth flow
type Oauth2Callback struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_url"`
}

// AccessTokenResponse contains metadata related to the access token
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// NewOauth2Application creates a new Oauth2 Application using the provided Client ID, Secret and scopes
func NewOauth2Application(clientID, clientSecret, scope string) *Oauth2Application {
	return &Oauth2Application{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}
}

// Callback generates an access_token from the supplied querystring parameters. Use this with the querystring parameters
// sent in the redirect url
func (o *Oauth2Application) Callback(obj Oauth2Callback) (*AccessTokenResponse, error) {
	parsed, _ := url.Parse("https://discordapp.com/api/oauth2/token")
	query := parsed.Query()
	query.Set("grant_type", "authorization_code")
	query.Set("code", obj.Code)
	query.Set("scope", o.Scope)
	query.Set("redirect_uri", obj.RedirectURI)
	parsed.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, parsed.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", basicAuth(o.ClientID, o.ClientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var resp *AccessTokenResponse
	json.Unmarshal(body, &resp)

	return resp, nil
}

// User returns the current authenticated user given the supplied access token
func (o *Oauth2Application) User(accessToken string) (u *User, err error) {
	// TODO: check scope
	req, err := http.NewRequest(http.MethodGet, "https://discordapp.com/api/v6/users/@me", nil)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}

	return
}

// Guilds returns an array of guilds the authenticated user is in
func (o *Oauth2Application) Guilds(accessToken string) (guilds []*Guild, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://discordapp.com/api/v6/users/@me/guilds", nil)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	err = json.Unmarshal(body, &guilds)
	if err != nil {
		return nil, err
	}

	return
}

// Implements basic HTTP authorization
func basicAuth(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
