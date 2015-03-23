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
//go:generate stringer -output types_standings_string.go -type=StandingType

package evego

// StandingType is the type of entity with which a Standing applies.
type StandingType int

// StandingType is the type of entity with which a Standing applies.
const (
	UnknownEntity StandingType = iota
	NPCFaction
	NPCCorporation
	NPCAgent
	PlayerCharacter
	PlayerCorporation
	PlayerAlliance
)

// Standing is a standing level from a specified entity.
type Standing struct {
	EntityType StandingType
	ID         int     `xml:"fromID,attr"`
	Name       string  `xml:"fromName,attr"`
	Standing   float64 `xml:"standing,attr"`
}
