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

package character_test

import (
	"database/sql"
	"testing"

	"github.com/backerman/evego/pkg/character"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStandingsCalculation(t *testing.T) {
	Convey("No standings = 0.0", t, func() {
		corp := sql.NullFloat64{Valid: false}
		fac := sql.NullFloat64{Valid: false}
		So(character.EffectiveStanding(corp, fac, 17, 666), ShouldEqual, 0.0)
	})

	Convey("An unskilled character has unchanged raw standing.", t, func() {
		corp := sql.NullFloat64{Valid: false}
		fac := sql.NullFloat64{Valid: true, Float64: 2.3}
		So(character.EffectiveStanding(corp, fac, 0, 0), ShouldEqual, 2.3)
	})

	Convey("Negative standings are affected by Diplomacy.", t, func() {
		corp := sql.NullFloat64{Valid: true, Float64: -3.0}
		fac := sql.NullFloat64{Valid: false}
		So(character.EffectiveStanding(corp, fac, 17, 3), ShouldEqual, -1.44)
	})

	Convey("Positive standings are affected by Connections.", t, func() {
		corp := sql.NullFloat64{Valid: true, Float64: 3.0}
		fac := sql.NullFloat64{Valid: false}
		So(character.EffectiveStanding(corp, fac, 3, -42), ShouldEqual, 3.84)
	})

	Convey("With corp and faction standings both present", t, func() {
		corp := sql.NullFloat64{Valid: true, Float64: 3.0}
		fac := sql.NullFloat64{Valid: true, Float64: -3.0}
		Convey("The higher of the two effective values is returned.", func() {
			So(character.EffectiveStanding(corp, fac, 3, 3), ShouldEqual, 3.84)
		})
	})

}
