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
//go:generate stringer -output items_string.go -type=ItemType

// Package types contains types that represent items in the EVE universe and
// are used by other packages in this library.
package types

import "fmt"

// ItemType is the type of an Eve item, and is used chiefly for
// correctly determining the item's reprocessing rate.
type ItemType int

// ItemType is the type of an Eve item.
const (
	UnknownItemType ItemType = iota
	Ore
	Ice
	Other
)

// Item is an Eve item.
type Item struct {
	Name      string `db:"typeName"`
	ID        int    `db:"typeID"`
	Type      ItemType
	Category  string `db:"categoryName"` // e.g. Module, Drone, Charge
	Group     string `db:"groupName"`    // e.g. Omber, Logistic Drone, Footwear
	BatchSize int    `db:"portionSize"`
}

func (i Item) String() string {
	return fmt.Sprintf("Item: %s (%d)", i.Name, i.ID)
}

// InventoryLine is an item in a material's composition, the player's
// inventory, or whatever.
type InventoryLine struct {
	Quantity int
	Item     *Item
}

func (i InventoryLine) String() string {
	return fmt.Sprintf("[%vx %v (%v)]", i.Quantity, i.Item.Name, i.Item.ID)
}

// MarketGroup is a group of items in the EVE market.
type MarketGroup struct {
	ID          int
	Parent      *MarketGroup
	Name        string
	Description string
}

func (m MarketGroup) String() string {
	result := fmt.Sprintf("%q (ID %d)", m.Name, m.ID)
	if m.Parent != nil {
		result += fmt.Sprintf(" parent: %v", m.Parent)
	} else {
		result += " parent nil"
	}
	return result
}
