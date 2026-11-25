/*
 *     Copyright 2026 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"

	"d7y.io/dragonfly/v2/scheduler/config"
)

// Cache is cache client.
type Cache struct {
	*cache.Cache
}

// New cache instance, it only uses the local cache if rdb is nil.
func New(rdb redis.UniversalClient) *Cache {
	// The attribute TTL of cache.Item(cache's instance) is the redis expiration time,
	// config.DefaultLFUCacheTTL is the expiration time of TinyLFU. TinyLFU does not
	// use the attribute TTL, so the value has to carry its own expiration time.
	return &Cache{
		Cache: cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(config.DefaultLFUCacheSize, config.DefaultLFUCacheTTL),
		}),
	}
}
