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

// Package routing provides services to calculate paths between points in EVE.
package routing

import (
	"io"

	"github.com/backerman/evego/pkg/types"
)

// EveRouter is the interface for the backing service that provides a router.
type EveRouter interface {
	io.Closer
	// NumJumps returns the number of jumps in the shortest path from
	// fromSystem to toSystem, or -1 if the destination is unreachable
	// from the start.
	NumJumps(fromSystem, toSystem *types.SolarSystem) (int, error)

	// NumJumpsID is a convenience method for NumJumps. Or is it the reverse?
	NumJumpsID(fromSystemID, toSystemID int) (int, error)
}
