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
	"time"

	"github.com/backerman/evego/pkg/types"
)

func unmarshalInt(cdata string, dest *int) {
	num, _ := strconv.ParseInt(string(cdata), 10, 0)
	*dest = int(num)
}

// Functions on a struct must be defined in the struct's package, so we can't
// define an unmarshaler for types.CharacterSheet here in package eveapi. The
// solution is to define a new type that's just types.CharacterSheet and then
// convert back before returning the types.CharacterSheet to the caller.

type charSheet types.CharacterSheet

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
					skillRow := types.Skill{}
					d.DecodeElement(&skillRow, &tok)
					cs.Skills = append(cs.Skills, skillRow)
				}
			}
			// if this is a row in a rowset, act based on parent rowset.
		case xml.EndElement:
			// do nothing
		case xml.CharData:
			// FIXME don't bother if all whitespace.
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

func (x *xmlAPI) CharacterSheet(characterID, keyID int, verificationCode string) (*types.CharacterSheet, time.Time, error) {
	params := url.Values{}
	params.Set("keyID", fmt.Sprintf("%d", keyID))
	params.Set("characterID", fmt.Sprintf("%d", characterID))
	params.Set("vcode", verificationCode)
	xmlBytes, err := x.get(characterSheet, params)
	if err != nil {
		return nil, time.Now(), err
	}
	var response charSheetAPIResponse
	xml.Unmarshal(xmlBytes, &response)
	// Convert back to types.CharacterSheet
	sheet := types.CharacterSheet(response.CharSheet)
	// Look up the name of each skill
	for i := range sheet.Skills {
		skill := &sheet.Skills[i]
		skillItem, err := x.db.ItemForID(skill.TypeID)
		if err != nil {
			skill.Name = fmt.Sprintf("Unknown skill (%d)", skill.TypeID)
		} else {
			skill.Name = skillItem.Name
		}
	}
	return &sheet, expirationTime(response.CurrentTime, response.CachedUntil), nil
}
