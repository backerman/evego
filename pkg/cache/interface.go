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

// Package cache defines an interface for caller-provided caches.
package cache

import (
	"io"
	"time"
)

// Cache is the interface expected by evego packages for a local cache.
type Cache interface {
	io.Closer

	// Get returns the cached value of the key, if it is available and unexpired.
	// It also returns the item's expiration time and a boolean flag that is true
	// if there was a hit, and false if the item was not in the cache or it was
	// expired.
	Get(key string) (*[]byte, time.Time, bool)

	// Put takes a key and blob to persist in the cache, and the item's expiry
	// time. It returns an error if something has gone wrong with the cache,
	// or nil otherwise.
	Put(key string, val *[]byte, expires time.Time) error
}
