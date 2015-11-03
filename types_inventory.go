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
//go:generate stringer -output types_inventory_string.go -type=InventoryFlag,BlueprintType

package evego

// InventoryItem is exactly what it sounds like.
type InventoryItem struct {
	// ItemID is a unique identifier for this object.
	ItemID int `xml:"itemID,attr"`
	// LocationID is the solar system or station where an object is located.
	LocationID int `xml:"locationID,attr"`
	// TypeID is the item's type.
	TypeID int `xml:"typeID,attr"`
	// Quantity is the number of items in this stack.
	Quantity int `xml:"quantity,attr"`
	// BlueprintType is the type of blueprint (original or copy), if applicable.
	BlueprintType BlueprintType `xml:"rawQuantity,attr"`
	// Unpackaged is true iff the item is unpackaged.
	Unpackaged bool `xml:"singleton,attr"`
	// Flag indicates the item's position; see the InventoryFlag enum.
	Flag InventoryFlag `xml:"flag,attr"`
	// Contents is a list of items that this item contains, if any.
	Contents []InventoryItem `xml:"rowset>row"`
}

// BlueprintItem is a blueprint returned from a blueprint endpoint.
type BlueprintItem struct {
	// ItemID is a unique identifier for this object.
	ItemID int `xml:"itemid,attr"`
	// LocationID is the solar system, station, or container ID where an object is located.
	LocationID int `xml:"locationid,attr"`
	// StationID is the solar system, station, or outpost ID where an object is located.
	StationID int `xml:"-"`
	// TypeID is the item's type.
	TypeID int `xml:"typeid,attr"`
	// TypeName is the item's type name.
	TypeName string `xml:"typename,attr"`
	// Quantity is the number of items in this stack, or -1 if it's not stacked.
	Quantity int `xml:"quantity,attr"`
	// Flag indicates the item's position; see the InventoryFlag enum.
	Flag InventoryFlag `xml:"flagid,attr"`
	// TimeEfficiency is the blueprint's researched time efficiency level [0..20]
	TimeEfficiency int `xml:"timeefficiency,attr"`
	// MaterialEfficiency is the blueprint's researched material efficiency level [0..10]
	MaterialEfficiency int `xml:"materialefficiency,attr"`
	// NumRuns is the number of runs remaining (-1 for original blueprints)
	NumRuns int `xml:"runs,attr"`
	// IsOriginal is true iff this blueprint is an original.
	IsOriginal bool `xml:"-"`
}

// BlueprintType is a blueprint's original/copy status.
type BlueprintType int

const (
	// NotBlueprint is not a blueprint.
	NotBlueprint BlueprintType = 0
	// BlueprintOriginal is an orignal blueprint.
	BlueprintOriginal BlueprintType = -1
	// BlueprintCopy is a blueprint copy.
	BlueprintCopy BlueprintType = -2
)

// InventoryFlag describes the location of an item in the asset list.
type InventoryFlag int

// This list of flags was dumped from the invFlags table.

const (
	// InvNone : None
	InvNone InventoryFlag = 0
	// InvWallet : Wallet
	InvWallet InventoryFlag = 1
	// InvFactory : Factory
	InvFactory InventoryFlag = 2
	// InvWardrobe : Wardrobe
	InvWardrobe InventoryFlag = 3
	// InvHangar : Hangar
	InvHangar InventoryFlag = 4
	// InvCargo : Cargo
	InvCargo InventoryFlag = 5
	// InvBriefcase : Briefcase
	InvBriefcase InventoryFlag = 6
	// InvSkill : Skill
	InvSkill InventoryFlag = 7
	// InvReward : Reward
	InvReward InventoryFlag = 8
	// InvConnected : Character in station connected
	InvConnected InventoryFlag = 9
	// InvDisconnected : Character in station offline
	InvDisconnected InventoryFlag = 10
	// InvLoSlot0 : Low power slot 1
	InvLoSlot0 InventoryFlag = 11
	// InvLoSlot1 : Low power slot 2
	InvLoSlot1 InventoryFlag = 12
	// InvLoSlot2 : Low power slot 3
	InvLoSlot2 InventoryFlag = 13
	// InvLoSlot3 : Low power slot 4
	InvLoSlot3 InventoryFlag = 14
	// InvLoSlot4 : Low power slot 5
	InvLoSlot4 InventoryFlag = 15
	// InvLoSlot5 : Low power slot 6
	InvLoSlot5 InventoryFlag = 16
	// InvLoSlot6 : Low power slot 7
	InvLoSlot6 InventoryFlag = 17
	// InvLoSlot7 : Low power slot 8
	InvLoSlot7 InventoryFlag = 18
	// InvMedSlot0 : Medium power slot 1
	InvMedSlot0 InventoryFlag = 19
	// InvMedSlot1 : Medium power slot 2
	InvMedSlot1 InventoryFlag = 20
	// InvMedSlot2 : Medium power slot 3
	InvMedSlot2 InventoryFlag = 21
	// InvMedSlot3 : Medium power slot 4
	InvMedSlot3 InventoryFlag = 22
	// InvMedSlot4 : Medium power slot 5
	InvMedSlot4 InventoryFlag = 23
	// InvMedSlot5 : Medium power slot 6
	InvMedSlot5 InventoryFlag = 24
	// InvMedSlot6 : Medium power slot 7
	InvMedSlot6 InventoryFlag = 25
	// InvMedSlot7 : Medium power slot 8
	InvMedSlot7 InventoryFlag = 26
	// InvHiSlot0 : High power slot 1
	InvHiSlot0 InventoryFlag = 27
	// InvHiSlot1 : High power slot 2
	InvHiSlot1 InventoryFlag = 28
	// InvHiSlot2 : High power slot 3
	InvHiSlot2 InventoryFlag = 29
	// InvHiSlot3 : High power slot 4
	InvHiSlot3 InventoryFlag = 30
	// InvHiSlot4 : High power slot 5
	InvHiSlot4 InventoryFlag = 31
	// InvHiSlot5 : High power slot 6
	InvHiSlot5 InventoryFlag = 32
	// InvHiSlot6 : High power slot 7
	InvHiSlot6 InventoryFlag = 33
	// InvHiSlot7 : High power slot 8
	InvHiSlot7 InventoryFlag = 34
	// InvFixedSlot : Fixed Slot
	InvFixedSlot InventoryFlag = 35
	// InvPromenadeSlot1 : Promenade Slot 1
	InvPromenadeSlot1 InventoryFlag = 40
	// InvPromenadeSlot2 : Promenade Slot 2
	InvPromenadeSlot2 InventoryFlag = 41
	// InvPromenadeSlot3 : Promenade Slot 3
	InvPromenadeSlot3 InventoryFlag = 42
	// InvPromenadeSlot4 : Promenade Slot 4
	InvPromenadeSlot4 InventoryFlag = 43
	// InvPromenadeSlot5 : Promenade Slot 5
	InvPromenadeSlot5 InventoryFlag = 44
	// InvPromenadeSlot6 : Promenade Slot 6
	InvPromenadeSlot6 InventoryFlag = 45
	// InvPromenadeSlot7 : Promenade Slot 7
	InvPromenadeSlot7 InventoryFlag = 46
	// InvPromenadeSlot8 : Promenade Slot 8
	InvPromenadeSlot8 InventoryFlag = 47
	// InvPromenadeSlot9 : Promenade Slot 9
	InvPromenadeSlot9 InventoryFlag = 48
	// InvPromenadeSlot10 : Promenade Slot 10
	InvPromenadeSlot10 InventoryFlag = 49
	// InvPromenadeSlot11 : Promenade Slot 11
	InvPromenadeSlot11 InventoryFlag = 50
	// InvPromenadeSlot12 : Promenade Slot 12
	InvPromenadeSlot12 InventoryFlag = 51
	// InvPromenadeSlot13 : Promenade Slot 13
	InvPromenadeSlot13 InventoryFlag = 52
	// InvPromenadeSlot14 : Promenade Slot 14
	InvPromenadeSlot14 InventoryFlag = 53
	// InvPromenadeSlot15 : Promenade Slot 15
	InvPromenadeSlot15 InventoryFlag = 54
	// InvPromenadeSlot16 : Promenade Slot 16
	InvPromenadeSlot16 InventoryFlag = 55
	// InvCapsule : Capsule
	InvCapsule InventoryFlag = 56
	// InvPilot : Pilot
	InvPilot InventoryFlag = 57
	// InvPassenger : Passenger
	InvPassenger InventoryFlag = 58
	// InvBoardingGate : Boarding gate
	InvBoardingGate InventoryFlag = 59
	// InvCrew : Crew
	InvCrew InventoryFlag = 60
	// InvSkillInTraining : Skill in training
	InvSkillInTraining InventoryFlag = 61
	// InvCorpMarket : Corporation Market Deliveries / Returns
	InvCorpMarket InventoryFlag = 62
	// InvLocked : Locked item, can not be moved unless unlocked
	InvLocked InventoryFlag = 63
	// InvUnlocked : Unlocked item, can be moved
	InvUnlocked InventoryFlag = 64
	// InvOfficeSlot1 : Office slot 1
	InvOfficeSlot1 InventoryFlag = 70
	// InvOfficeSlot2 : Office slot 2
	InvOfficeSlot2 InventoryFlag = 71
	// InvOfficeSlot3 : Office slot 3
	InvOfficeSlot3 InventoryFlag = 72
	// InvOfficeSlot4 : Office slot 4
	InvOfficeSlot4 InventoryFlag = 73
	// InvOfficeSlot5 : Office slot 5
	InvOfficeSlot5 InventoryFlag = 74
	// InvOfficeSlot6 : Office slot 6
	InvOfficeSlot6 InventoryFlag = 75
	// InvOfficeSlot7 : Office slot 7
	InvOfficeSlot7 InventoryFlag = 76
	// InvOfficeSlot8 : Office slot 8
	InvOfficeSlot8 InventoryFlag = 77
	// InvOfficeSlot9 : Office slot 9
	InvOfficeSlot9 InventoryFlag = 78
	// InvOfficeSlot10 : Office slot 10
	InvOfficeSlot10 InventoryFlag = 79
	// InvOfficeSlot11 : Office slot 11
	InvOfficeSlot11 InventoryFlag = 80
	// InvOfficeSlot12 : Office slot 12
	InvOfficeSlot12 InventoryFlag = 81
	// InvOfficeSlot13 : Office slot 13
	InvOfficeSlot13 InventoryFlag = 82
	// InvOfficeSlot14 : Office slot 14
	InvOfficeSlot14 InventoryFlag = 83
	// InvOfficeSlot15 : Office slot 15
	InvOfficeSlot15 InventoryFlag = 84
	// InvOfficeSlot16 : Office slot 16
	InvOfficeSlot16 InventoryFlag = 85
	// InvBonus : Bonus
	InvBonus InventoryFlag = 86
	// InvDroneBay : Drone Bay
	InvDroneBay InventoryFlag = 87
	// InvBooster : Booster
	InvBooster InventoryFlag = 88
	// InvImplant : Implant
	InvImplant InventoryFlag = 89
	// InvShipHangar : Ship Hangar
	InvShipHangar InventoryFlag = 90
	// InvShipOffline : Ship Offline
	InvShipOffline InventoryFlag = 91
	// InvRigSlot0 : Rig power slot 1
	InvRigSlot0 InventoryFlag = 92
	// InvRigSlot1 : Rig power slot 2
	InvRigSlot1 InventoryFlag = 93
	// InvRigSlot2 : Rig power slot 3
	InvRigSlot2 InventoryFlag = 94
	// InvRigSlot3 : Rig power slot 4
	InvRigSlot3 InventoryFlag = 95
	// InvRigSlot4 : Rig power slot 5
	InvRigSlot4 InventoryFlag = 96
	// InvRigSlot5 : Rig power slot 6
	InvRigSlot5 InventoryFlag = 97
	// InvRigSlot6 : Rig power slot 7
	InvRigSlot6 InventoryFlag = 98
	// InvRigSlot7 : Rig power slot 8
	InvRigSlot7 InventoryFlag = 99
	// InvFactoryOperation : Factory Background Operation
	InvFactoryOperation InventoryFlag = 100
	// InvCorpSAG2 : Corp Security Access Group 2
	InvCorpSAG2 InventoryFlag = 116
	// InvCorpSAG3 : Corp Security Access Group 3
	InvCorpSAG3 InventoryFlag = 117
	// InvCorpSAG4 : Corp Security Access Group 4
	InvCorpSAG4 InventoryFlag = 118
	// InvCorpSAG5 : Corp Security Access Group 5
	InvCorpSAG5 InventoryFlag = 119
	// InvCorpSAG6 : Corp Security Access Group 6
	InvCorpSAG6 InventoryFlag = 120
	// InvCorpSAG7 : Corp Security Access Group 7
	InvCorpSAG7 InventoryFlag = 121
	// InvSecondaryStorage : Secondary Storage
	InvSecondaryStorage InventoryFlag = 122
	// InvCaptainsQuarters : Captains Quarters
	InvCaptainsQuarters InventoryFlag = 123
	// InvWisPromenade : Wis Promenade
	InvWisPromenade InventoryFlag = 124
	// InvSubSystem0 : Sub system slot 0
	InvSubSystem0 InventoryFlag = 125
	// InvSubSystem1 : Sub system slot 1
	InvSubSystem1 InventoryFlag = 126
	// InvSubSystem2 : Sub system slot 2
	InvSubSystem2 InventoryFlag = 127
	// InvSubSystem3 : Sub system slot 3
	InvSubSystem3 InventoryFlag = 128
	// InvSubSystem4 : Sub system slot 4
	InvSubSystem4 InventoryFlag = 129
	// InvSubSystem5 : Sub system slot 5
	InvSubSystem5 InventoryFlag = 130
	// InvSubSystem6 : Sub system slot 6
	InvSubSystem6 InventoryFlag = 131
	// InvSubSystem7 : Sub system slot 7
	InvSubSystem7 InventoryFlag = 132
	// InvSpecializedFuelBay : Specialized Fuel Bay
	InvSpecializedFuelBay InventoryFlag = 133
	// InvSpecializedOreHold : Specialized Ore Hold
	InvSpecializedOreHold InventoryFlag = 134
	// InvSpecializedGasHold : Specialized Gas Hold
	InvSpecializedGasHold InventoryFlag = 135
	// InvSpecializedMineralHold : Specialized Mineral Hold
	InvSpecializedMineralHold InventoryFlag = 136
	// InvSpecializedSalvageHold : Specialized Salvage Hold
	InvSpecializedSalvageHold InventoryFlag = 137
	// InvSpecializedShipHold : Specialized Ship Hold
	InvSpecializedShipHold InventoryFlag = 138
	// InvSpecializedSmallShipHold : Specialized Small Ship Hold
	InvSpecializedSmallShipHold InventoryFlag = 139
	// InvSpecializedMediumShipHold : Specialized Medium Ship Hold
	InvSpecializedMediumShipHold InventoryFlag = 140
	// InvSpecializedLargeShipHold : Specialized Large Ship Hold
	InvSpecializedLargeShipHold InventoryFlag = 141
	// InvSpecializedIndustrialShipHold : Specialized Industrial Ship Hold
	InvSpecializedIndustrialShipHold InventoryFlag = 142
	// InvSpecializedAmmoHold : Specialized Ammo Hold
	InvSpecializedAmmoHold InventoryFlag = 143
	// InvStructureActive : StructureActive
	InvStructureActive InventoryFlag = 144
	// InvStructureInactive : StructureInactive
	InvStructureInactive InventoryFlag = 145
	// InvJunkyardReprocessed : This item was put into a junkyard through reprocession.
	InvJunkyardReprocessed InventoryFlag = 146
	// InvJunkyardTrashed : This item was put into a junkyard through being trashed by its owner.
	InvJunkyardTrashed InventoryFlag = 147
	// InvSpecializedCommandCenterHold : Specialized Command Center Hold
	InvSpecializedCommandCenterHold InventoryFlag = 148
	// InvSpecializedPlanetaryCommoditiesHold : Specialized Planetary Commodities Hold
	InvSpecializedPlanetaryCommoditiesHold InventoryFlag = 149
	// InvPlanetSurface : Planet Surface
	InvPlanetSurface InventoryFlag = 150
	// InvSpecializedMaterialBay : Specialized Material Bay
	InvSpecializedMaterialBay InventoryFlag = 151
	// InvDustCharacterDatabank : Dust Character Databank
	InvDustCharacterDatabank InventoryFlag = 152
	// InvDustCharacterBattle : Dust Character Battle
	InvDustCharacterBattle InventoryFlag = 153
	// InvQuafeBay : Quafe Bay
	InvQuafeBay InventoryFlag = 154
	// InvFleetHangar : Fleet Hangar
	InvFleetHangar InventoryFlag = 155
	// InvHiddenModifiers : Hidden Modifiers
	InvHiddenModifiers InventoryFlag = 156
)
