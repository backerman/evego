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
	"io"
	"time"

	"github.com/backerman/evego/pkg/types"
)

// XMLAPI is an interface to the EVE XML API. We could make the interface
// sufficiently abstract to cover multiple APIs, but that seems on the silly
// side.
type XMLAPI interface {
	io.Closer

	// OutpostForID returns a conquerable station with the provided ID.
	OutpostForID(id int) (*types.Station, error)

	// OutpostsForName returns the stations matching the provided name pattern.
	// The percent character (%) may be used as a wildcard.
	OutpostsForName(name string) ([]types.Station, error)

	// AccountCharacters returns a list of characters that the provided key can
	// access. It also returns the expiration time for this information; the caller
	// must cache the returned data until that time.
	AccountCharacters(key *XMLKey) ([]types.Character, time.Time, error)

	// CharacterSheet returns the character sheet for the given character ID.
	// It also returns the expiration time for this information; the caller must
	// cache the returned data until that time.
	CharacterSheet(key *XMLKey, characterID int) (*types.CharacterSheet, time.Time, error)

	// CharacterStandings returns a character's standings. It also returns the
	// expiration time for this information; the callermust cache the returned
	// data until that time.
	CharacterStandings(key *XMLKey, characterID int) ([]types.Standing, time.Time, error)
}
