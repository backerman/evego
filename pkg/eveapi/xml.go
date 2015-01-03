/*
Copyright © 2014 Brad Ackerman.

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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/types"
)

var (
	cacheExpiry = time.Time{}
	outposts    = make(map[int]*types.Station)
)

const (
	iso8601            = "2006-01-02 15:04:05"
	conqerableStations = "/eve/ConquerableStationList.xml.aspx"
)

type xmlAPI struct {
	// Endpoint URL to access.
	url  *url.URL
	http http.Client
	db   dbaccess.EveDatabase
}

// XMLAPI returns an EveAPI that accesses the EVE Online XML API.
func XMLAPI(serviceURL string, staticDB dbaccess.EveDatabase) EveAPI {
	endpoint, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatalf("Unable to process endpoint URL: %v", err)
	}
	return &xmlAPI{url: endpoint, db: staticDB}
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

func (x *xmlAPI) get(endpoint string, params ...url.Values) ([]byte, error) {
	// Make a copy of our base URL and modify as appropriate for this call.
	callURL := *x.url
	callURL.Path = endpoint
	var (
		req *http.Request
		err error
	)
	if params == nil {
		req, err = http.NewRequest("GET", callURL.String(), nil)
	} else {
		req, err = http.NewRequest("POST", callURL.String(),
			strings.NewReader(params[0].Encode()))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "evego (https://github.com/backerman/evego)")
	if params != nil {
		// We're doing a POST so need to set the content type.
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := x.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (x *xmlAPI) Close() error {
	return nil
}
