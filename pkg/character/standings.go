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

// Package character provides calculation of character statistics.
package character

import (
	"database/sql"
	"math"
)

// Stolen from https://gist.github.com/DavidVaini/10308388
func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

// EffectiveStanding calculates a character's effective standing towards
// an NPC corporation.
func EffectiveStanding(rawCorp, rawFaction sql.NullFloat64, connections, diplomacy int) float64 {
	effective := func(raw float64) float64 {
		// Apply the standings equation.
		var skill int
		if raw < 0 {
			skill = diplomacy
		} else {
			skill = connections
		}
		return 10.0 - (10.0-raw)*(1-0.04*float64(skill))
	}
	effStandings := make([]float64, 0, 2)
	// Check each standing; append to effStandings.
	for _, s := range []sql.NullFloat64{rawCorp, rawFaction} {
		if s.Valid {
			effStandings = append(effStandings, effective(s.Float64))
		}
	}
	// Pick the maximum standing and return it.
	switch len(effStandings) {
	case 0:
		// Neither standing was present; default is zero.
		return 0
	case 1:
		// Only one of corp/faction standings exists; return it.
		return round(effStandings[0], 2)
	default:
		// Toon has both corp and faction standing - return whichever is greater.
		return round(math.Max(effStandings[0], effStandings[1]), 2)
	}
}
