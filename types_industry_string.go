// generated by stringer -output types_industry_string.go -type=ActivityType; DO NOT EDIT

package evego

import "fmt"

const _ActivityType_name = "NoneManufacturingResearchingTechnologyResearchingTEResearchingMECopyingDuplicatingReverseEngineeringInvention"

var _ActivityType_index = [...]uint8{0, 4, 17, 38, 51, 64, 71, 82, 100, 109}

func (i ActivityType) String() string {
	if i < 0 || i+1 >= ActivityType(len(_ActivityType_index)) {
		return fmt.Sprintf("ActivityType(%d)", i)
	}
	return _ActivityType_name[_ActivityType_index[i]:_ActivityType_index[i+1]]
}
