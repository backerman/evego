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

package cache_test

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/backerman/evego/pkg/cache"
	"github.com/garyburd/redigo/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testRedis = ":6379" // localhost, default port
)

func TestRedisCache(t *testing.T) {

	Convey("Set up redis", t, func() {
		pool := &redis.Pool{
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", testRedis)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}

		testCache := cache.RedisCache(testRedis)

		Convey("Putting a thing to the cache succeeds.", func() {
			key := "Testing! " + strconv.FormatInt(time.Now().UnixNano(), 10)
			value, expiry := "Fred", time.Now().Add(15*time.Second)
			conn := pool.Get()
			defer conn.Close()

			// Shouldn't be there until it's created
			resp, err := conn.Do("GET", key)
			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)

			// Now we create it.
			err = testCache.Put(key, []byte(value), expiry)
			So(err, ShouldBeNil)

			// Retrieve from Redis instance and verify value.
			stored, err := redis.Bytes(conn.Do("GET", key))
			So(err, ShouldBeNil)

			// stored should be value, but zlib-compressed.
			compressedBuffer := bytes.NewBuffer(stored)
			r, err := zlib.NewReader(compressedBuffer)
			So(err, ShouldBeNil)
			readFromRedis, err := ioutil.ReadAll(r)
			So(err, ShouldBeNil)
			r.Close()
			So(readFromRedis, ShouldResemble, []byte(value))

			Convey("Getting the thing should also succeed.", func() {
				// Pull from the cache; verify that it exists and was retrieved correctly.
				fromCache, found := testCache.Get(key)
				So(found, ShouldBeTrue)
				So(fromCache, ShouldResemble, []byte(value))

				Convey("TTL of cached values should be correctly set.", func() {
					ttl, err := redis.Int(conn.Do("TTL", key))
					So(err, ShouldBeNil)
					So(ttl, ShouldBeBetween, 10, 15)
				})
			})
		})

		Convey("Getting something not in the cache fails.", func() {
			key := "This had better not be there. " + strconv.FormatInt(time.Now().UnixNano(), 10)
			_, found := testCache.Get(key)
			So(found, ShouldBeFalse)
		})

		Convey("Closing the cache should not fail.", func() {
			err := testCache.Close()
			So(err, ShouldBeNil)
		})
	})
}
