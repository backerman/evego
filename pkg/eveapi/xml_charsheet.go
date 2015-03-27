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
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/backerman/evego"
)

func unmarshalInt(cdata string, dest *int) {
	num, _ := strconv.ParseInt(string(cdata), 10, 0)
	*dest = int(num)
}

// Functions on a struct must be defined in the struct's package, so we can't
// define an unmarshaler for evego.CharacterSheet here in package eveapi. The
// solution is to define a new type that's just evego.CharacterSheet and then
// convert back before returning the evego.CharacterSheet to the caller.

type charSheet evego.CharacterSheet

func (cs *charSheet) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var currentElement, rowsetName string
	for {
		token, err := d.Token()
		switch err {
		case nil:
		// pass
		case io.EOF:
			return nil
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
				switch rowsetName {
				case "skills":
					skillRow := evego.Skill{}
					d.DecodeElement(&skillRow, &tok)
					cs.Skills = append(cs.Skills, skillRow)
				}
			}
			// if this is a row in a rowset, act based on parent rowset.
		case xml.EndElement:
			// do nothing
		case xml.CharData:
			contents := strings.TrimSpace(string(tok))
			if len(contents) > 0 {
				switch currentElement {
				case "characterID":
					unmarshalInt(contents, &cs.ID)
				case "name":
					cs.Name = string(contents)
				case "corporationName":
					cs.Corporation = string(contents)
				case "corporationID":
					unmarshalInt(contents, &cs.CorporationID)
				case "allianceName":
					cs.Alliance = string(contents)
				case "allianceID":
					unmarshalInt(contents, &cs.AllianceID)
				}
			}
		}
	}
}

type charSheetAPIResponse struct {
	CurrentTime string    `xml:"currentTime"`
	CharSheet   charSheet `xml:"result"`
	CachedUntil string    `xml:"cachedUntil"`
}

func (x *xmlAPI) CharacterSheet(key *evego.XMLKey, characterID int) (*evego.CharacterSheet, error) {
	params := url.Values{}
	params.Set("keyID", strconv.Itoa(key.KeyID))
	params.Set("characterID", strconv.Itoa(characterID))
	params.Set("vcode", key.VerificationCode)
	xmlBytes, err := x.get(characterSheet, params)
	if err != nil {
		return nil, err
	}
	var response charSheetAPIResponse
	xml.Unmarshal(xmlBytes, &response)
	// Convert back to evego.CharacterSheet
	sheet := evego.CharacterSheet(response.CharSheet)
	// Look up the name of each skill
	for i := range sheet.Skills {
		skill := &sheet.Skills[i]
		skillItem, err := x.db.ItemForID(skill.TypeID)
		if err != nil {
			skill.Name = fmt.Sprintf("Unknown skill (%d)", skill.TypeID)
			skill.Group = fmt.Sprintf("Unknown")
		} else {
			skill.Name = skillItem.Name
			skill.Group = skillItem.Group
		}
	}
	return &sheet, nil
}
