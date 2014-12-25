// generated by stringer -output industry_string.go -type=ActivityType; DO NOT EDIT

package types

import "fmt"

const _ActivityType_name = "NoneManufacturingResearchingTechnologyResearchingTEResearchingMECopyingDuplicatingReverseEngineeringInvention"

var _ActivityType_index = [...]uint8{4, 17, 38, 51, 64, 71, 82, 100, 109}

func (i ActivityType) String() string {
	if i < 0 || i >= ActivityType(len(_ActivityType_index)) {
		return fmt.Sprintf("ActivityType(%d)", i)
	}
	hi := _ActivityType_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _ActivityType_index[i-1]
	}
	return _ActivityType_name[lo:hi]
}
