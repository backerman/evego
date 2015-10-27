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
	"database/sql"
	"encoding/xml"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/backerman/evego"
)

// XML headers for unmarshalling

type outpostAPIResponse struct {
	CurrentTime string    `xml:"currentTime"`
	Outposts    []outpost `xml:"result>rowset>row"`
	CachedUntil string    `xml:"cachedUntil"`
}

type outpost struct {
	Name            string `xml:"stationName,attr"`
	ID              int    `xml:"stationID,attr"`
	TypeID          int    `xml:"stationTypeID,attr"`
	SolarSystemID   int    `xml:"solarSystemID,attr"`
	CorporationID   int    `xml:"corporationID,attr"`
	CorporationName string `xml:"corporationName,attr"`
}

// Check the cache expiry time and update it if necessary.
func (x *xmlAPI) checkOutpostCache() error {
	if time.Now().After(cacheExpiry) {
		// The cache has expired or has not yet been populated.
		// FIXME Only one goroutine should update cache. Mutex? Channel?
		newOutposts := make(map[int]*evego.Station)
		xmlBytes, err := x.get(conqerableStations)
		if err != nil {
			// FIXME some sort of throttling required
			return err
		}
		var response outpostAPIResponse
		xml.Unmarshal(xmlBytes, &response)
		cacheExpiry = expirationTime(response.CurrentTime, response.CachedUntil)
		for i := range response.Outposts {
			o := response.Outposts[i]
			stn := evego.Station{
				Name:          o.Name,
				ID:            o.ID,
				SystemID:      o.SolarSystemID,
				Corporation:   o.CorporationName,
				CorporationID: o.CorporationID,
			}
			system, err := x.db.SolarSystemForID(stn.SystemID)
			if err != nil {
				return err
			}
			stn.ConstellationID = system.ConstellationID
			stn.RegionID = system.RegionID
			newOutposts[o.ID] = &stn
		}
		outposts = newOutposts
	}

	return nil
}

func (x *xmlAPI) OutpostForID(id int) (*evego.Station, error) {
	err := x.checkOutpostCache()
	if err != nil {
		return nil, err
	}
	stn, exists := outposts[id]
	if !exists {
		return nil, fmt.Errorf("Station ID %d not found.", id)
	}
	return stn, nil
}

// Implementation of sort.Interface

type stationsList []evego.Station

func (s stationsList) Len() int {
	return len(s)
}

func (s stationsList) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s stationsList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (x *xmlAPI) OutpostsForName(name string) ([]evego.Station, error) {
	// This is a horribly inefficient implementation. Switch to SQLite
	// in-memory DB rather than keeping everything as Golang structs?
	err := x.checkOutpostCache()
	if err != nil {
		return nil, err
	}
	var stations []evego.Station
	namePattern := "^(?i:" + strings.Replace(name, "%", ".*", -1) + ")$"
	nameRE, err := regexp.Compile(namePattern)
	if err != nil {
		return nil, err
	}
	for id, stn := range outposts {
		if nameRE.MatchString(stn.Name) {
			matchStn, err := x.OutpostForID(id)
			if err != nil {
				return nil, err
			}
			stations = append(stations, *matchStn)
		}
	}

	if len(stations) == 0 {
		return stations, sql.ErrNoRows
	}
	sort.Sort(stationsList(stations))
	return stations, nil
}

func (x *xmlAPI) DumpOutposts() []*evego.Station {
	// Refresh if necessary.
	x.checkOutpostCache()
	// ... and dump 'em.
	outpostsSlice := make([]*evego.Station, 0, len(outposts))
	for _, outpost := range outposts {
		outpostsSlice = append(outpostsSlice, outpost)
	}
	return outpostsSlice
}
