// generated by stringer -output types_inventory_string.go -type=InventoryFlag,BlueprintType; DO NOT EDIT

package evego

import "fmt"

const (
	_InventoryFlag_name_0 = "InvNoneInvWalletInvFactoryInvWardrobeInvHangarInvCargoInvBriefcaseInvSkillInvRewardInvConnectedInvDisconnectedInvLoSlot0InvLoSlot1InvLoSlot2InvLoSlot3InvLoSlot4InvLoSlot5InvLoSlot6InvLoSlot7InvMedSlot0InvMedSlot1InvMedSlot2InvMedSlot3InvMedSlot4InvMedSlot5InvMedSlot6InvMedSlot7InvHiSlot0InvHiSlot1InvHiSlot2InvHiSlot3InvHiSlot4InvHiSlot5InvHiSlot6InvHiSlot7InvFixedSlot"
	_InventoryFlag_name_1 = "InvPromenadeSlot1InvPromenadeSlot2InvPromenadeSlot3InvPromenadeSlot4InvPromenadeSlot5InvPromenadeSlot6InvPromenadeSlot7InvPromenadeSlot8InvPromenadeSlot9InvPromenadeSlot10InvPromenadeSlot11InvPromenadeSlot12InvPromenadeSlot13InvPromenadeSlot14InvPromenadeSlot15InvPromenadeSlot16InvCapsuleInvPilotInvPassengerInvBoardingGateInvCrewInvSkillInTrainingInvCorpMarketInvLockedInvUnlocked"
	_InventoryFlag_name_2 = "InvOfficeSlot1InvOfficeSlot2InvOfficeSlot3InvOfficeSlot4InvOfficeSlot5InvOfficeSlot6InvOfficeSlot7InvOfficeSlot8InvOfficeSlot9InvOfficeSlot10InvOfficeSlot11InvOfficeSlot12InvOfficeSlot13InvOfficeSlot14InvOfficeSlot15InvOfficeSlot16InvBonusInvDroneBayInvBoosterInvImplantInvShipHangarInvShipOfflineInvRigSlot0InvRigSlot1InvRigSlot2InvRigSlot3InvRigSlot4InvRigSlot5InvRigSlot6InvRigSlot7InvFactoryOperation"
	_InventoryFlag_name_3 = "InvCorpSAG2InvCorpSAG3InvCorpSAG4InvCorpSAG5InvCorpSAG6InvCorpSAG7InvSecondaryStorageInvCaptainsQuartersInvWisPromenadeInvSubSystem0InvSubSystem1InvSubSystem2InvSubSystem3InvSubSystem4InvSubSystem5InvSubSystem6InvSubSystem7InvSpecializedFuelBayInvSpecializedOreHoldInvSpecializedGasHoldInvSpecializedMineralHoldInvSpecializedSalvageHoldInvSpecializedShipHoldInvSpecializedSmallShipHoldInvSpecializedMediumShipHoldInvSpecializedLargeShipHoldInvSpecializedIndustrialShipHoldInvSpecializedAmmoHoldInvStructureActiveInvStructureInactiveInvJunkyardReprocessedInvJunkyardTrashedInvSpecializedCommandCenterHoldInvSpecializedPlanetaryCommoditiesHoldInvPlanetSurfaceInvSpecializedMaterialBayInvDustCharacterDatabankInvDustCharacterBattleInvQuafeBayInvFleetHangarInvHiddenModifiers"
)

var (
	_InventoryFlag_index_0 = [...]uint16{0, 7, 16, 26, 37, 46, 54, 66, 74, 83, 95, 110, 120, 130, 140, 150, 160, 170, 180, 190, 201, 212, 223, 234, 245, 256, 267, 278, 288, 298, 308, 318, 328, 338, 348, 358, 370}
	_InventoryFlag_index_1 = [...]uint16{0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 171, 189, 207, 225, 243, 261, 279, 289, 297, 309, 324, 331, 349, 362, 371, 382}
	_InventoryFlag_index_2 = [...]uint16{0, 14, 28, 42, 56, 70, 84, 98, 112, 126, 141, 156, 171, 186, 201, 216, 231, 239, 250, 260, 270, 283, 297, 308, 319, 330, 341, 352, 363, 374, 385, 404}
	_InventoryFlag_index_3 = [...]uint16{0, 11, 22, 33, 44, 55, 66, 85, 104, 119, 132, 145, 158, 171, 184, 197, 210, 223, 244, 265, 286, 311, 336, 358, 385, 413, 440, 472, 494, 512, 532, 554, 572, 603, 641, 657, 682, 706, 728, 739, 753, 771}
)

func (i InventoryFlag) String() string {
	switch {
	case 0 <= i && i <= 35:
		return _InventoryFlag_name_0[_InventoryFlag_index_0[i]:_InventoryFlag_index_0[i+1]]
	case 40 <= i && i <= 64:
		i -= 40
		return _InventoryFlag_name_1[_InventoryFlag_index_1[i]:_InventoryFlag_index_1[i+1]]
	case 70 <= i && i <= 100:
		i -= 70
		return _InventoryFlag_name_2[_InventoryFlag_index_2[i]:_InventoryFlag_index_2[i+1]]
	case 116 <= i && i <= 156:
		i -= 116
		return _InventoryFlag_name_3[_InventoryFlag_index_3[i]:_InventoryFlag_index_3[i+1]]
	default:
		return fmt.Sprintf("InventoryFlag(%d)", i)
	}
}

const _BlueprintType_name = "BlueprintCopyBlueprintOriginalNotBlueprint"

var _BlueprintType_index = [...]uint8{0, 13, 30, 42}

func (i BlueprintType) String() string {
	i -= -2
	if i < 0 || i >= BlueprintType(len(_BlueprintType_index)-1) {
		return fmt.Sprintf("BlueprintType(%d)", i+-2)
	}
	return _BlueprintType_name[_BlueprintType_index[i]:_BlueprintType_index[i+1]]
}