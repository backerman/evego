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
//go:generate stringer -output industry_string.go -type=ActivityType

package types

import "fmt"

// ActivityType is an industrial activity performed on or resulting in
// a blueprint.
type ActivityType int

// The ActivityType values.
const (
	None ActivityType = iota
	Manufacturing
	ResearchingTechnology
	ResearchingTE
	ResearchingME
	Copying
	Duplicating
	ReverseEngineering
	Invention
)

// IndustryActivity is an action (e.g. invention) taken on an input item
// (e.g. Vexor Blueprint) producing a result (e.g. Ishtar Blueprint).
type IndustryActivity struct {
	InputItem      *Item
	ActivityType   ActivityType
	OutputItem     *Item
	OutputQuantity int
}

func (i IndustryActivity) String() string {
	return fmt.Sprintf("Activity %v: %v -> %d x %v", i.ActivityType, i.InputItem,
		i.OutputQuantity, i.OutputItem)
}
