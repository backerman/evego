// generated by stringer -output standings_string.go -type=StandingType; DO NOT EDIT

package types

import "fmt"

const _StandingType_name = "UnknownEntityNPCFactionNPCCorporationNPCAgentPlayerCharacterPlayerCorporationPlayerAlliance"

var _StandingType_index = [...]uint8{13, 23, 37, 45, 60, 77, 91}

func (i StandingType) String() string {
	if i < 0 || i >= StandingType(len(_StandingType_index)) {
		return fmt.Sprintf("StandingType(%d)", i)
	}
	hi := _StandingType_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _StandingType_index[i-1]
	}
	return _StandingType_name[lo:hi]
}
