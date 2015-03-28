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

package evego_test

import (
	"testing"

	. "github.com/backerman/evego"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSkillSort(t *testing.T) {
	Convey("Given an array of skills", t, func() {
		skillArray := []Skill{
			{Name: "Gunnery", Group: "Gunnery", TypeID: 3300, NumSkillpoints: 256000, Level: 5, Published: true},
			{Name: "Small Hybrid Turret", Group: "Gunnery", TypeID: 3301, NumSkillpoints: 256000, Level: 5, Published: true},
			{Name: "Spaceship Command", Group: "Spaceship Command", TypeID: 3327, NumSkillpoints: 45255, Level: 4, Published: true},
			{Name: "Gallente Frigate", Group: "Spaceship Command", TypeID: 3328, NumSkillpoints: 512000, Level: 5, Published: true},
			{Name: "Mining", Group: "Resource Processing", TypeID: 3386, NumSkillpoints: 256000, Level: 5, Published: true},
		}

		Convey("They are sorted correctly.", func() {
			expected := []Skill{
				{Name: "Gunnery", Group: "Gunnery", TypeID: 3300, NumSkillpoints: 256000, Level: 5, Published: true},
				{Name: "Small Hybrid Turret", Group: "Gunnery", TypeID: 3301, NumSkillpoints: 256000, Level: 5, Published: true},
				{Name: "Mining", Group: "Resource Processing", TypeID: 3386, NumSkillpoints: 256000, Level: 5, Published: true},
				{Name: "Gallente Frigate", Group: "Spaceship Command", TypeID: 3328, NumSkillpoints: 512000, Level: 5, Published: true},
				{Name: "Spaceship Command", Group: "Spaceship Command", TypeID: 3327, NumSkillpoints: 45255, Level: 4, Published: true},
			}
			SortSkills(skillArray)
			So(skillArray, ShouldResemble, expected)
		})
	})
}
