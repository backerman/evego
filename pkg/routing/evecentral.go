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

// Package routing provides services to calculate paths between points in EVE.
package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/backerman/evego"
)

// Yes, the http.Client here is nil. We don't actually need to set anything
// on it; all we use it for is to get access to the Do method.

type eveCentralRouter struct {
	http      http.Client
	endpoint  *url.URL
	respCache evego.Cache
}

// jsonResponse is the response provided by the EVE-Central server.
type jsonResponse struct {
	From      jsonSystem `json:"from"`
	To        jsonSystem `json:"to"`
	SecChange bool       `json:"secChange"`
}

type jsonSystem struct {
	ID       int     `json:"systemID"`
	Name     string  `json:"name"`
	Security float64 `json:"security"`
}

// EveCentralRouter creates an EveRouter that uses EVE-Central's API
// to provide routing. The endpoint argument will normally be
// "http://api.eve-central.com/api/route".
func EveCentralRouter(endpoint string, aCache evego.Cache) evego.Router {
	epURL, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("Invalid URL %v passed for Eve-Central endpoint: %v", endpoint, err)
	}
	return &eveCentralRouter{endpoint: epURL, respCache: aCache}
}

func (r *eveCentralRouter) getURL(u string) ([]byte, error) {
	// Check cache first.
	cachedBody, found := r.respCache.Get(u)
	if found {
		return cachedBody, nil
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "evego (https://github.com/backerman/evego)")
	resp, err := r.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	// EVE-Central doesn't specify a caching time to use, so we're picking
	// five minutes at random.
	r.respCache.Put(u, body, time.Now().Add(5*time.Minute))
	return body, err
}

func (r *eveCentralRouter) NumJumps(fromSystem, toSystem *evego.SolarSystem) (int, error) {
	return r.NumJumpsID(fromSystem.ID, toSystem.ID)
}

func (r *eveCentralRouter) NumJumpsID(fromSystemID, toSystemID int) (int, error) {
	// Don't even query the server if the start and end are identical.
	if fromSystemID == toSystemID {
		return 0, nil
	}
	// Copy the endpoint.
	queryURL := *r.endpoint
	queryURL.Path += fmt.Sprintf("/from/%d/to/%d", fromSystemID, toSystemID)
	respJSON, err := r.getURL(queryURL.String())
	if err != nil {
		return 0, err
	}
	var resp []jsonResponse
	json.Unmarshal(respJSON, &resp)
	if len(resp) == 0 {
		// Since we know the start and end are not identical, this means that you
		// can't get there from here; in this case, the specification says that
		// we return -1.
		return -1, nil
	}

	// Otherwise, resp has one element for each jump.
	return len(resp), nil
}

func (r *eveCentralRouter) Close() error {
	return nil
}
