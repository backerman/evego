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
	"net/url"
	"strconv"

	"github.com/backerman/evego"
)

const (
	blueprintStr = "Blueprint"
	blueprintLen = len(blueprintStr)
)

type assetsList []evego.InventoryItem

type assetsResponse struct {
	CurrentTime string     `xml:"currentTime"`
	Assets      assetsList `xml:"result>rowset>row"`
	CachedUntil string     `xml:"cachedUntil"`
}

func (x *xmlAPI) processAssets(assets []evego.InventoryItem) error {
	for i := range assets {
		asset := &assets[i]
		thisAsset, err := x.db.ItemForID(asset.TypeID)
		if err != nil {
			return err
		}
		startIndex := len(thisAsset.Name) - blueprintLen
		var endOfName string
		if startIndex > 0 {
			endOfName = thisAsset.Name[startIndex:]
		}
		if endOfName != "Blueprint" {
			asset.BlueprintType = evego.NotBlueprint
		} else if !asset.Unpackaged {
			// This is a blueprint, but it's packaged, and therefore cannot be
			// a copy.
			asset.BlueprintType = evego.BlueprintOriginal
		}
		if asset.Quantity == 0 {
			// The default quantity is 1, and our default is 0, so fix.
			asset.Quantity = 1
		}
		if asset.Contents == nil {
			// No contents, but we want to make sure there's a slice here (rather than
			// just a nil) for consistency.
			asset.Contents = make([]evego.InventoryItem, 0, 0)
		} else if len(asset.Contents) > 0 {
			err = x.processAssets(asset.Contents)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (x *xmlAPI) Assets(key *evego.XMLKey, characterID int) ([]evego.InventoryItem, error) {
	params := url.Values{}
	params.Set("keyID", strconv.Itoa(key.KeyID))
	params.Set("characterID", strconv.Itoa(characterID))
	params.Set("vcode", key.VerificationCode)
	xmlBytes, err := x.get(characterAssets, params)
	if err != nil {
		return nil, err
	}
	var response assetsResponse
	xml.Unmarshal(xmlBytes, &response)
	assets := []evego.InventoryItem(response.Assets)
	err = x.processAssets(assets)
	return assets, err
}
