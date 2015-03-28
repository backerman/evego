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

package evego

import "sort"

// Character represents one EVE player toon.
type Character struct {
	Name          string `json:"name"          xml:"name,attr"`
	ID            int    `json:"id"            xml:"characterID,attr"`
	Corporation   string `json:"corporation"   xml:"corporationName,attr"`
	CorporationID int    `json:"corporationID" xml:"corporationID,attr"`
	Alliance      string `json:"alliance"      xml:"allianceName,attr"`
	AllianceID    int    `json:"allianceID"    xml:"allianceID,attr"`
}

// CharacterSheet contains the character sheet information for a toon
// as provied by the /char/CharacterSheet.xml.aspx endpoint.
type CharacterSheet struct {
	Character
	Skills []Skill `json:"skills"`
}

// Skill is one of a character's injected skills.
type Skill struct {
	Name           string `json:"name"`
	Group          string `json:"group"`
	TypeID         int    `json:"typeID"         xml:"typeID,attr"`
	NumSkillpoints int    `json:"numSkillpoints" xml:"skillpoints,attr"`
	Level          int    `json:"level"          xml:"level,attr"`
	Published      bool   `json:"isPublished"    xml:"published,attr"`
}

// Wrappers to sort skills using the standard library's sort package.
type skillsSorted []Skill

func (s skillsSorted) Len() int      { return len(s) }
func (s skillsSorted) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s skillsSorted) Less(i, j int) bool {
	return s[i].Group < s[j].Group || (s[i].Group == s[j].Group && s[i].Name < s[j].Name)
}

// SortSkills sorts the provided skill array by group and name.
func SortSkills(skills []Skill) {
	sortMe := skillsSorted(skills)
	sort.Sort(sortMe)
}
