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

// Package eveapi is the public interface for accessing the
// EVE APIs (XML, CREST, or whatever.)
package eveapi

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/backerman/evego"
)

var (
	cacheExpiry = time.Time{}
	outposts    = make(map[int]*evego.Station)
)

const (
	iso8601            = "2006-01-02 15:04:05"
	accountCharacters  = "/account/Characters.xml.aspx"
	characterAssets    = "/char/AssetList.xml.aspx"
	characterSheet     = "/char/CharacterSheet.xml.aspx"
	characterStandings = "/char/Standings.xml.aspx"
	conqerableStations = "/eve/ConquerableStationList.xml.aspx"
)

type xmlAPI struct {
	// Endpoint URL to access.
	url   *url.URL
	http  http.Client
	db    evego.Database
	cache evego.Cache
}

// XML returns an object that accesses the EVE Online XML API.
func XML(serviceURL string, staticDB evego.Database, aCache evego.Cache) evego.XMLAPI {
	endpoint, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatalf("Unable to process endpoint URL: %v", err)
	}
	return &xmlAPI{url: endpoint, db: staticDB, cache: aCache}
}

func expirationTime(currentTime, cachedUntil string) time.Time {
	// This system's time *should* be identical to the server's,
	// but we're not taking any chances.
	current, err := time.Parse(iso8601, currentTime)
	if err != nil {
		// If we can't parse, return the datum.
		log.Printf("Unable to parse time from server: %v", currentTime)
		return time.Time{}
	}
	until, err := time.Parse(iso8601, cachedUntil)
	if err != nil {
		log.Printf("Unable to parse time from server: %v", cachedUntil)
		return time.Time{}
	}
	diff := until.Sub(current)
	return time.Now().Add(diff)
}

// get will do an extra unmarshal - it just wants the current time and
// expiry time.
type expiryInfo struct {
	CurrentTime string `xml:"currentTime"`
	CachedUntil string `xml:"cachedUntil"`
}

// get executes a call to the EVE API given the endpoint path and an optional
// url.Values containing the parameters to be passed.
func (x *xmlAPI) get(endpoint string, params ...url.Values) ([]byte, error) {
	// Make a copy of our base URL and modify as appropriate for this call.
	callURL := *x.url
	callURL.Path = endpoint
	var (
		req *http.Request
		err error
	)
	if params != nil {
		callURL.RawQuery = params[0].Encode()
	}

	urlStr := callURL.String()
	// Check cache.
	cachedBody, found := x.cache.Get(urlStr)
	if found {
		return cachedBody, nil
	}

	req, err = http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "evego (https://github.com/backerman/evego)")
	resp, err := x.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Put our repsonse in the cache.
	expiry := expiryInfo{}
	xml.Unmarshal(body, &expiry)
	expiresAt := expirationTime(expiry.CurrentTime, expiry.CachedUntil)
	x.cache.Put(urlStr, body, expiresAt)
	return body, err
}

func (x *xmlAPI) Close() error {
	return nil
}
