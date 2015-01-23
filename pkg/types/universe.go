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

package types

// SolarSystem is a solar system within the EVE universe.
type SolarSystem struct {
	Name            string `db:"solarSystemName"`
	ID              int    `db:"solarSystemID"`
	Constellation   string `db:"constellationName"`
	ConstellationID int    `db:"constellationID"`
	Region          string `db:"regionName"`
	RegionID        int    `db:"regionID"`
	Security        float64
}

// Region is one of the regions in the EVE universe.
type Region struct {
	Name string `db:"regionName"`
	ID   int    `db:"regionID"`
}

// Station is either an NPC station or a conquerable outpost.
type Station struct {
	Name                   string  `db:"stationName"`
	ID                     int     `db:"stationID"`
	SystemID               int     `db:"solarSystemID"`
	ConstellationID        int     `db:"constellationID"`
	RegionID               int     `db:"regionID"`
	CorporationID          int     `db:"corporationID"`
	Corporation            string  `db:"corporationName"`
	ReprocessingEfficiency float64 `db:"reprocessingEfficiency"`
}
