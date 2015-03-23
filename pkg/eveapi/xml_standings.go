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

package eveapi

import (
	"encoding/xml"
	"io"
	"log"
	"net/url"
	"strconv"

	"github.com/backerman/evego"
)

// Functions on a struct must be defined in the struct's package, so we can't
// define an unmarshaler for evego.CharacterSheet here in package eveapi. The
// solution is to define a new type that's just evego.CharacterSheet and then
// convert back before returning the evego.CharacterSheet to the caller.

type standingsList []evego.Standing

type standingsResponse struct {
	CurrentTime string        `xml:"currentTime"`
	Standings   standingsList `xml:"result"`
	CachedUntil string        `xml:"cachedUntil"`
}

func (sl *standingsList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var currentElement, rowsetName string
	standings := make([]evego.Standing, 0, 2)
	done := false
	for !done {
		token, err := d.Token()
		switch err {
		case nil:
		// pass
		case io.EOF:
			done = true
		default:
			log.Fatalf("This should not happen: %v", err)
		}
		switch tok := token.(type) {
		case xml.StartElement:
			currentElement = tok.Name.Local
			switch currentElement {
			case "rowset":
				for _, a := range tok.Attr {
					if a.Name.Local == "name" {
						rowsetName = a.Value
						break
					}
				}
			case "row":
				standing := evego.Standing{}
				d.DecodeElement(&standing, &tok)
				switch rowsetName {
				case "agents":
					standing.EntityType = evego.NPCAgent
				case "NPCCorporations":
					standing.EntityType = evego.NPCCorporation
				case "factions":
					standing.EntityType = evego.NPCFaction
				}
				standings = append(standings, standing)
			}
			// if this is a row in a rowset, act based on parent rowset.
		case xml.EndElement:
			// do nothing
		case xml.CharData:
			// do nothing
		}
	}

	*sl = standings
	return nil
}

func (x *xmlAPI) CharacterStandings(key *evego.XMLKey, characterID int) ([]evego.Standing, error) {
	params := url.Values{}
	params.Set("keyID", strconv.Itoa(key.KeyID))
	params.Set("characterID", strconv.Itoa(characterID))
	params.Set("vcode", key.VerificationCode)
	xmlBytes, err := x.get(characterStandings, params)
	if err != nil {
		return nil, err
	}
	var response standingsResponse
	xml.Unmarshal(xmlBytes, &response)
	standings := []evego.Standing(response.Standings)
	return standings, nil
}
