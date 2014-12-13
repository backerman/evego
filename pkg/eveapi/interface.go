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

// Package eveapi is the public interface for accessing the
// EVE APIs (XML, CREST, or whatever.)
package eveapi

import (
	"io"

	"github.com/backerman/evego/pkg/types"
)

// EveAPI is an interface to the EVE API.
type EveAPI interface {
	io.Closer

	// OutpostForID returns a conquerable station with the provided ID.
	OutpostForID(id int) (*types.Station, error)
}
