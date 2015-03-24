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
//go:generate stringer -output types_market_string.go -type=OrderType,OrderRange

package evego

import (
	"fmt"
	"time"
)

// OrderType is the order's type (either buy or sell); or All for searching
// for either
type OrderType int

const (
	// Buy order
	Buy OrderType = iota
	// Sell order
	Sell
	// AllOrders is either buy or sell (used for searches only)
	AllOrders
)

// OrderRange is the area from which capsuleers can sell to
// a buy order.
type OrderRange int

const (
	// BuyStation is only the order's station.
	BuyStation OrderRange = iota
	// BuySystem is any station in the order's system.
	BuySystem
	// BuyNumberJumps is a specified number of jumps from the order's system.
	BuyNumberJumps
	// BuyRegion is anywhere within the order's region.
	BuyRegion
)

// Order represents an order on the EVE market.
type Order struct {
	Type       OrderType
	Item       *Item
	Quantity   int
	Price      float64
	Station    *Station
	Expiration time.Time

	// Fields for buy orders only

	MinQuantity int
	JumpRange   OrderRange
	NumJumps    int
}

func (o *Order) String() string {
	var out string
	switch o.Type {
	case Sell:
		out += "Sell order"
	case Buy:
		out += "Buy order"
	default:
		out += "Order of invalid type"
	}
	out += fmt.Sprintf(" for %d units of", o.Quantity)
	if o.Item != nil {
		out += " " + o.Item.Name
	} else {
		out += " an invalid item"
	}
	out += fmt.Sprintf(" at %.2f ISK each in %v", o.Price, o.Station.Name)

	return out
}
