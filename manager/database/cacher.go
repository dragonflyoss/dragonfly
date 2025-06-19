/*
 *     Copyright 2024 The Dragonfly Authors
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

package database

import (
	"context"
	"time"

	xxhash "github.com/cespare/xxhash/v2"
	lru "github.com/elastic/go-freelru"
	caches "github.com/go-gorm/caches/v4"
)

const (
	// defaultCacheSize is the default size of the cache.
	defaultCacheSize = 1024

	// defaultTTL is the default TTL of the cache.
	defaultTTL = 30 * time.Second
)

// cacher is a cache implementation using LRU for gorm.
type cacher struct {
	// store is the LRU cache.
	store *lru.ShardedLRU[string, any]
}

// hashStringXXHASH returns the hash of the string s.
func hashStringXXHASH(s string) uint32 {
	return uint32(xxhash.Sum64String(s))
}

// newCacher creates a new cacher.
func newCacher() (caches.Cacher, error) {
	store, err := lru.NewSharded[string, any](defaultCacheSize, hashStringXXHASH)
	if err != nil {
		return nil, err
	}

	store.SetLifetime(defaultTTL)
	return &cacher{store: store}, nil
}

// Get impl should check if a specific key exists in the cache and return its value
// look at Query.Marshal.
func (c *cacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	val, ok := c.store.Get(key)
	if !ok {
		return nil, nil
	}

	if err := q.Unmarshal(val.([]byte)); err != nil {
		return nil, err
	}

	return q, nil
}

// Store impl should store a cached representation of the val param
// look at Query.Unmarshal.
func (c *cacher) Store(ctx context.Context, key string, val *caches.Query[any]) error {
	res, err := val.Marshal()
	if err != nil {
		return err
	}

	c.store.Add(key, res)
	return nil
}

// Invalidate impl should invalidate all cached values. It will be called when
// INSERT / UPDATE / DELETE queries are sent to the DB.
func (c *cacher) Invalidate(ctx context.Context) error {
	var err error
	c.store, err = lru.NewSharded[string, any](defaultCacheSize, hashStringXXHASH)
	if err != nil {
		return err
	}

	c.store.SetLifetime(defaultTTL)
	return nil
}
