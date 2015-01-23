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

// Character represents one EVE player toon.
type Character struct {
	Name          string `json:"name"          xml:"name,attr"`
	ID            int    `json:"id"            xml:"characterID,attr"`
	Corporation   string `json:"corporation"   xml:"corporationName,attr"`
	CorporationID int    `json:"corporationID" xml:"corporationID,attr"`
	Alliance      string `json:"alliance"      xml:"allianceName,attr"`
	AllianceID    int    `json:"allianceID"    xml:"allianceID,attr"`
}

// CharacterSheet ...
type CharacterSheet struct {
	Name          string  `json:"name"`
	ID            int     `json:"id"`
	Corporation   string  `json:"corporation"`
	CorporationID int     `json:"corporationID"`
	Alliance      string  `json:"alliance"`
	AllianceID    int     `json:"allianceID"`
	Skills        []Skill `json:"skills"`
}

// Skill ...
type Skill struct {
	Name           string `json:"name"`
	TypeID         int    `json:"typeID"         xml:"typeID,attr"`
	NumSkillpoints int    `json:"numSkillpoints" xml:"skillpoints,attr"`
	Level          int    `json:"level"          xml:"level,attr"`
	Published      bool   `json:"isPublished"    xml:"published,attr"`
}
