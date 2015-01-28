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

package cache

import "time"

type nilCache struct {
}

// NilCache returns a struct satisfying the Cache interface that goes nowhere
// and does nothing. All calls will GNDN, and Get will always result in a cache
// miss.
func NilCache() Cache {
	return &nilCache{}
}

func (c *nilCache) Close() error {
	return nil
}

func (c *nilCache) Get(key string) (*[]byte, time.Time, bool) {
	return nil, time.Time{}, false
}

func (c *nilCache) Put(key string, val *[]byte, expires time.Time) error {
	return nil
}