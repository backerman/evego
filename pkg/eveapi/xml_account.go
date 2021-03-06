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
	"net/url"

	"github.com/backerman/evego"
)

type charsResponse struct {
	CurrentTime string            `xml:"currentTime"`
	Characters  []evego.Character `xml:"result>rowset>row"`
	CachedUntil string            `xml:"cachedUntil"`
}

func (x *xmlAPI) AccountCharacters(key *evego.XMLKey) ([]evego.Character, error) {
	params := url.Values{}
	params.Set("keyID", fmt.Sprintf("%d", key.KeyID))
	params.Set("vcode", key.VerificationCode)
	xmlBytes, err := x.get(accountCharacters, params)
	if err != nil {
		return nil, err
	}
	var response charsResponse
	xml.Unmarshal(xmlBytes, &response)

	return response.Characters, nil
}
