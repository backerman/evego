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

// Package market is the public interface for accessing
// EVE market data.
package market

import (
	"io"

	"github.com/backerman/evego/pkg/types"
)

// EveMarket returns information about market orders.
type EveMarket interface {
	io.Closer

	// OrdersForItem returns the market orders for a given item.
	// location is the name of either a system or a region.
	// type can be Buy, Sell, or All.
	OrdersForItem(itemID *types.Item, location string, orderType types.OrderType) (*[]types.Order, error)
}
