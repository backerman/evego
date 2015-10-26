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

package test

import (
	"fmt"
	"time"

	"github.com/backerman/evego"
)

// CacheData maintains information on calls by the test code to the cache.
type CacheData struct {
	GetKeys, PutKeys map[string]bool
	NumGets, NumPuts int
	PutExpires       time.Time
	// When ShouldError is set to true, a call to Put will error out.
	ShouldError bool
}

type testCache struct {
	data *CacheData
}

// Cache returns a cache object used for testing.
func Cache(data *CacheData) evego.Cache {
	data.GetKeys = make(map[string]bool)
	data.PutKeys = make(map[string]bool)
	return &testCache{
		data: data,
	}
}

func (c *testCache) Get(key string) ([]byte, bool) {
	c.data.GetKeys[key] = true
	c.data.NumGets++
	return nil, false
}

func (c *testCache) Put(key string, val []byte, expires time.Time) error {
	c.data.PutKeys[key] = true
	c.data.NumPuts++
	c.data.PutExpires = expires
	if c.data.ShouldError {
		return fmt.Errorf("Error requested")
	}
	return nil
}

func (c *testCache) Close() error {
	return nil
}
