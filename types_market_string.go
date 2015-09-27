// generated by stringer -output types_market_string.go -type=OrderType,OrderRange; DO NOT EDIT

package evego

import "fmt"

const _OrderType_name = "BuySellAllOrders"

var _OrderType_index = [...]uint8{0, 3, 7, 16}

func (i OrderType) String() string {
	if i < 0 || i >= OrderType(len(_OrderType_index)-1) {
		return fmt.Sprintf("OrderType(%d)", i)
	}
	return _OrderType_name[_OrderType_index[i]:_OrderType_index[i+1]]
}

const _OrderRange_name = "BuyStationBuySystemBuyNumberJumpsBuyRegion"

var _OrderRange_index = [...]uint8{0, 10, 19, 33, 42}

func (i OrderRange) String() string {
	if i < 0 || i >= OrderRange(len(_OrderRange_index)-1) {
		return fmt.Sprintf("OrderRange(%d)", i)
	}
	return _OrderRange_name[_OrderRange_index[i]:_OrderRange_index[i+1]]
}
