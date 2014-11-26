/*
Copyright © 2014 Brad Ackerman.

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

// Package industry calculates the output from manufacturing, reprocessing, and
// other industrial processes.
package industry

import (
	"math"
	"strings"

	"github.com/backerman/evego/pkg/types"
)

// ReproSkills is a character's skills that are applicable to reprocessing.
type ReproSkills struct {
	Reprocessing           int
	ReprocessingEfficiency int
	ScrapmetalProcessing   int
	OreProcessing          map[string]int
}

// StationTax returns the tax rate for a station based on the character's
// standing.
func StationTax(standing float64) float64 {
	taxRate := math.Max(0.0, 0.05-0.0075*standing)
	return taxRate
}

// round returns the nearest integer to the input: floor(in) if the fractional
// component of in is less than 0.5, and ceil(in) otherwise.
func round(in float64) float64 {
	i, f := math.Modf(in)
	if f < 0.5 {
		return i
	}
	return i + 1.0
}

// ReprocessItem returns the result of reprocessing a given item and the number
// of input items that were reprocessed.
func reprocessItem(item *types.Item, quantity int, stationYield float64, standing float64, skills ReproSkills) *[]types.InventoryLine {
	yield := stationYield
	switch item.Type {
	case types.Ice, types.Ore:
		splitName := strings.Split(item.Name, " ")
		baseName := splitName[len(splitName)-1]
		yield *= 1.0 + float64(skills.Reprocessing)*0.03
		yield *= 1.0 + float64(skills.ReprocessingEfficiency)*0.02
		yield *= 1.0 + float64(skills.OreProcessing[baseName])*0.02
	default:
		yield *= 1.0 + float64(skills.ScrapmetalProcessing)*0.02
	}
	reprocessed := []types.InventoryLine{}

	// Ensure that the quantity is okay.
	batch := item.BatchSize
	quantity, remainder := quantity/batch, quantity%batch
	if remainder != 0 {
		// When the quantity of an item is not an integer multiple of its
		// batch size, pass through the fractional batch unprocessed.
		reprocessed = append(reprocessed,
			types.InventoryLine{Quantity: remainder, Item: item})
	}

	// Add yielded items
	for _, el := range item.Materials {
		// Take station tax based on the truncated number of units produced.
		newQuantity := math.Floor(float64(quantity*el.Quantity) * yield)
		stationCut := round(newQuantity * StationTax(standing))
		newQuantity = newQuantity - stationCut
		quantInt := int(newQuantity)
		if quantInt > 0 {
			reprocessed = append(reprocessed,
				types.InventoryLine{Quantity: quantInt, Item: el.Item})
		}

	}
	return &reprocessed
}

// ReprocessItems reprocesses a number of items, consolidating stacks of each
// output item.
func ReprocessItems(items *[]types.InventoryLine, stationYield float64, standing float64, skills ReproSkills) *[]types.InventoryLine {

	reproed := []types.InventoryLine{}
	for _, item := range *items {
		outItems := reprocessItem(item.Item, item.Quantity, stationYield, standing, skills)
		reproed = append(reproed, *outItems...)
	}

	// Deduplicate items
	quantities := make(map[int]int)
	outItems := make(map[int]*types.Item)
	for _, line := range reproed {
		quantities[line.Item.ID] += line.Quantity
		outItems[line.Item.ID] = line.Item
	}
	// blank return array
	reproed = []types.InventoryLine{}
	for _, item := range outItems {
		reproed = append(reproed,
			types.InventoryLine{Quantity: quantities[item.ID], Item: outItems[item.ID]})
	}
	return &reproed
}

// ReprocessItem is a convenience function for reprocessing a single item.
func ReprocessItem(item *types.Item, quantity int, stationYield float64, standing float64, skills ReproSkills) *[]types.InventoryLine {
	items := &[]types.InventoryLine{
		{Item: item, Quantity: quantity},
	}

	return ReprocessItems(items, stationYield, standing, skills)
}
