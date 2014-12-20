/*
Copyright © 2014 Brad Ackerman.

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

// Package dbaccess is the public interface for accessing the
// EVE static data export.
package dbaccess

import (
	"io"

	"github.com/backerman/evego/pkg/types"
)

// EveDatabase is an object that returns information about items in EVE.
type EveDatabase interface {
	io.Closer

	// Items

	ItemForName(itemName string) (*types.Item, error)
	MarketGroupForItem(item *types.Item) (*types.MarketGroup, error)

	// Universe locations

	SolarSystemForID(systemID int) (*types.SolarSystem, error)
	SolarSystemForName(systemName string) (*types.SolarSystem, error)
	RegionForName(regionName string) (*types.Region, error)
	StationForID(stationID int) (*types.Station, error)

	// Blueprints, invention, and manufacturing

	// BlueprintOutputs returns the items and quantity of each that can be output
	// by performing industrial actions on a blueprint given that blueprint's name
	// (typeName) as a string. The type name may include the percent (%) character
	// as a wildcard.
	BlueprintOutputs(typeName string) (*[]types.IndustryActivity, error)

	// BlueprintForProduct returns the blueprints that can produce a given output.
	BlueprintForProduct(typeName string) (*[]types.IndustryActivity, error)

	// BlueprintsUsingMaterial returns the blueprints that use the given input material
	// in an industrial process (manufacturing, invention, etc.)
	BlueprintsUsingMaterial(typeName string) (*[]types.IndustryActivity, error)

	// BlueprintProductionInputs returns the required materials for one run
	// of production on an unresearched (ME 0% / TE 0%) blueprint. It takes as
	// parameters the blueprint to be used and the selected output product.
	BlueprintProductionInputs(
		typeName string, outputTypeName string) (*[]types.InventoryLine, error)
}
