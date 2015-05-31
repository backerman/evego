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

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"log"
	"time"

	"github.com/backerman/evego"
	"github.com/garyburd/redigo/redis"
)

type redisCache struct {
	pool redis.Pool
}

// RedisCache returns a new evego.Cache object that uses a Redis backend.
// It takes a server hostname/port in the format accepted by net.Dial and,
// optionally, a password as parameters.
func RedisCache(server string, password ...string) evego.Cache {
	c := redisCache{
		pool: redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				if len(password) > 0 {
					// A password has been passed, so use it.
					if _, err := c.Do("AUTH", password[0]); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
	return &c
}

func (c *redisCache) Close() error {
	return c.pool.Close()
}

func (c *redisCache) Get(key string) ([]byte, bool) {
	conn := c.pool.Get()
	defer conn.Close()
	cachedBytes, err := redis.Bytes(conn.Do("GET", key))
	if err == nil {
		r, err := zlib.NewReader(bytes.NewReader(cachedBytes))
		if err != nil {
			log.Fatalf("Unable to decompress cached data: %v", err)
			return nil, false
		}
		gotten, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatalf("Unable to decompress cached data: %v", err)
			return nil, false
		}
		r.Close()
		return gotten, true
	} else if err == redis.ErrNil {
		return nil, false
	}
	// Some error other than not having found the cached response.
	log.Fatalf("Unable to access cache: %v", err)
	return nil, false
}

func (c *redisCache) Put(key string, val []byte, expires time.Time) error {
	conn := c.pool.Get()
	defer conn.Close()

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(val)
	w.Close()
	expiresSeconds := expires.Sub(time.Now()).Seconds()
	_, err := conn.Do("SET", key, b.Bytes(), "EX", int(expiresSeconds))
	return err
}
