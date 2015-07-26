/*
Copyright © 2014–5 Brad Ackerman.

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

// Package evesso provides support for EVE's single sign-on API.
package evesso

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"golang.org/x/oauth2"
)

var (
	// Endpoint is the production EVE Online cluster. (Tranquility)
	Endpoint = oauth2.Endpoint{
		AuthURL:  "https://login.eveonline.com/oauth/authorize",
		TokenURL: "https://login.eveonline.com/oauth/token",
	}

	// TestEndpoint is the public testing EVE Online cluster. (Singularity)
	TestEndpoint = oauth2.Endpoint{
		AuthURL:  "https://sisilogin.testeveonline.com/oauth/authorize",
		TokenURL: "https://sisilogin.testeveonline.com/oauth/token",
	}
)

const (
	// Scopes that can be requested.

	// PublicData is the only scope available at present.
	PublicData = "publicData"
)

// Authenticator provides an interface to EVE's OAuth2 interface.
type Authenticator interface {
	// URL returns the authentication URL for a given state value.
	URL(state string) string
	// Exchange converts an authorization code into an OAuth2 token.
	Exchange(code string) (*oauth2.Token, error)
	// CharacterInfo returns information on the authentication token's scopes
	// and the character that was authenticated.
	CharacterInfo(token *oauth2.Token) (*CharacterInfo, error)
}

type ssoAuthenticator struct {
	config oauth2.Config
}

// CharacterInfo is the information returned by the CREST API about the character
// that authenticated.
type CharacterInfo struct {
	CharacterID        int
	CharacterName      string
	Scopes             string
	TokenType          string
	CharacterOwnerHash string
}

// MakeAuthenticator returns an Authenticator object for the given client and
// access point. The provided client configuration (secret, ID, and redirect URL)
// must be registered for anything to work.
//
// Optionally, add one or more scopes to request.
func MakeAuthenticator(endpoint oauth2.Endpoint, clientID, clientSecret, redirectURL string, scopes ...string) Authenticator {
	a := ssoAuthenticator{}
	a.config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}
	return &a
}

func (a *ssoAuthenticator) URL(state string) string {
	return a.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (a *ssoAuthenticator) Exchange(code string) (*oauth2.Token, error) {
	return a.config.Exchange(oauth2.NoContext, code)
}

func (a *ssoAuthenticator) CharacterInfo(token *oauth2.Token) (*CharacterInfo, error) {
	authenticatedClient := a.config.Client(oauth2.NoContext, token)
	verifyURL := strings.Replace(a.config.Endpoint.AuthURL, "authorize", "verify", 1)
	resp, err := authenticatedClient.Get(verifyURL)
	if err != nil {
		return nil, err
	}
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	charInfo := &CharacterInfo{}
	err = json.Unmarshal(respJSON, charInfo)
	if err != nil {
		return nil, err
	}
	return charInfo, nil
}
