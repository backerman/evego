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

import (
	"io"
)

// Market returns information about market orders.
type Market interface {
	io.Closer

	// OrdersForItem returns the market orders for a given item.
	// location is the name of either a system or a region.
	// type can be Buy, Sell, or All.
	OrdersForItem(itemID *Item, location string, orderType OrderType) (*[]Order, error)

	// BuyInStation returns the buy orders that are in range of the given
	// station (i.e., can be sold to by a user there).
	BuyInStation(itemID *Item, location *Station) (*[]Order, error)

	// OrdersInStation returns the buy orders that are in range of a given station,
	// and the sell orders available at that station.
	OrdersInStation(item *Item, location *Station) (*[]Order, error)
}
