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

import "io"

// XMLKey is a key ID / verification code pair used to retrieve data from the
// EVE XML API.
type XMLKey struct {
	KeyID            int
	VerificationCode string
	// Description is an optional description provided by the user.
	Description string
}

// XMLAPI is an interface to the EVE XML API. We could make the interface
// sufficiently abstract to cover multiple APIs, but that seems on the silly
// side.
type XMLAPI interface {
	io.Closer

	// OutpostForID returns a conquerable station with the provided ID.
	OutpostForID(id int) (*Station, error)

	// OutpostsForName returns the stations matching the provided name pattern.
	// The percent character (%) may be used as a wildcard.
	OutpostsForName(name string) ([]Station, error)

	// DumpOutposts returns the current list of outposts.
	DumpOutposts() []*Station

	// AccountCharacters returns a list of characters that the provided key can
	// access.
	AccountCharacters(key *XMLKey) ([]Character, error)

	// CharacterSheet returns the character sheet for the given character ID.
	CharacterSheet(key *XMLKey, characterID int) (*CharacterSheet, error)

	// CharacterStandings returns a character's standings.
	CharacterStandings(key *XMLKey, characterID int) ([]Standing, error)

	// Assets gets a character's assets.
	Assets(key *XMLKey, characterID int) ([]InventoryItem, error)

	// Blueprints gets a character's blueprints.
	Blueprints(key *XMLKey, characterID int, assets []InventoryItem) ([]BlueprintItem, error)
}
