/*
 *     Copyright 2020 The Dragonfly Authors
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
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	v1 = "foo"
	v2 = "bar"
	v3 = "baz"
	v4 = "yes"
)

type mockStruct struct {
	Num      int
	Children []*mockStruct
}

func mockExpiredItem(c Cache, k string, x any) {
	c.(*cache).items[k] = Item{Object: x, Expiration: time.Now().Add(-time.Hour).UnixNano()}
}

func TestItem_Expired(t *testing.T) {
	tests := []struct {
		name   string
		item   Item
		expect func(t *testing.T, expired bool)
	}{
		{
			name: "zero expiration never expires",
			item: Item{Object: 1, Expiration: 0},
			expect: func(t *testing.T, expired bool) {
				assert := assert.New(t)
				assert.False(expired)
			},
		},
		{
			name: "future deadline is not expired",
			item: Item{Object: 1, Expiration: time.Now().Add(time.Hour).UnixNano()},
			expect: func(t *testing.T, expired bool) {
				assert := assert.New(t)
				assert.False(expired)
			},
		},
		{
			name: "past deadline is expired",
			item: Item{Object: 1, Expiration: time.Now().Add(-time.Hour).UnixNano()},
			expect: func(t *testing.T, expired bool) {
				assert := assert.New(t)
				assert.True(expired)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.item.Expired())
		})
	}
}

func TestCache(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		expect func(t *testing.T, x any, found bool)
	}{
		{
			name: "int value",
			key:  "a",
			expect: func(t *testing.T, x any, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal(1, x)
			},
		},
		{
			name: "string value",
			key:  "b",
			expect: func(t *testing.T, x any, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal("b", x)
			},
		},
		{
			name: "float value",
			key:  "c",
			expect: func(t *testing.T, x any, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal(3.5, x)
			},
		},
		{
			name: "missing key",
			key:  "d",
			expect: func(t *testing.T, x any, found bool) {
				assert := assert.New(t)
				assert.False(found)
				assert.Nil(x)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(DefaultExpiration, 0)
			c.Set("a", 1, DefaultExpiration)
			c.Set("b", "b", DefaultExpiration)
			c.Set("c", 3.5, DefaultExpiration)

			x, found := c.Get(tc.key)
			tc.expect(t, x, found)
		})
	}
}

func TestCacheTimes(t *testing.T) {
	assert := assert.New(t)
	c := New(100*time.Millisecond, 1*time.Millisecond)
	c.Set("a", 1, DefaultExpiration)
	c.Set("b", 2, NoExpiration)
	c.Set("c", 3, 40*time.Millisecond)
	c.Set("d", 4, 140*time.Millisecond)

	<-time.After(50 * time.Millisecond)
	_, found := c.Get("c")
	assert.False(found)

	<-time.After(80 * time.Millisecond)
	_, found = c.Get("a")
	assert.False(found)
	_, found = c.Get("b")
	assert.True(found)
	_, found = c.Get("d")
	assert.True(found)

	<-time.After(40 * time.Millisecond)
	_, found = c.Get("d")
	assert.False(found)
}

func TestStorePointerToStruct(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	foo := &mockStruct{Num: 1}
	c.Set(v1, foo, DefaultExpiration)

	x, found := c.Get(v1)
	assert.True(found)
	assert.Same(foo, x)

	foo.Num++
	y, found := c.Get(v1)
	assert.True(found)
	assert.Equal(&mockStruct{Num: 2}, y)
}

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		n       int
		expect  func(t *testing.T, keys []string, err error)
	}{
		{
			name:    "limit below number of matches",
			pattern: "^b",
			n:       1,
			expect: func(t *testing.T, keys []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(keys, 1)
			},
		},
		{
			name:    "limit equals number of matches",
			pattern: "^b",
			n:       2,
			expect: func(t *testing.T, keys []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(keys, 2)
			},
		},
		{
			name:    "limit above number of matches",
			pattern: "^b",
			n:       4,
			expect: func(t *testing.T, keys []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(keys, 2)
			},
		},
		{
			name:    "longer prefix matches both keys",
			pattern: "^ba",
			n:       2,
			expect: func(t *testing.T, keys []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(keys, 2)
			},
		},
		{
			name:    "pattern matches no key",
			pattern: "^a",
			n:       2,
			expect: func(t *testing.T, keys []string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(keys)
			},
		},
		{
			name:    "invalid regular expression",
			pattern: "(",
			n:       2,
			expect: func(t *testing.T, keys []string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(keys)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(DefaultExpiration, 0)
			c.Set(v2, v1, DefaultExpiration)
			c.Set(v3, v1, DefaultExpiration)

			keys, err := c.Scan(tc.pattern, tc.n)
			tc.expect(t, keys, err)
		})
	}
}

func TestSetDefault(t *testing.T) {
	tests := []struct {
		name              string
		defaultExpiration time.Duration
		expect            func(t *testing.T, expiration time.Time, found bool)
	}{
		{
			name:              "default expiration sets the deadline",
			defaultExpiration: time.Hour,
			expect: func(t *testing.T, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.WithinDuration(time.Now().Add(time.Hour), expiration, time.Minute)
			},
		},
		{
			name:              "no expiration leaves the item without deadline",
			defaultExpiration: NoExpiration,
			expect: func(t *testing.T, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.True(expiration.IsZero())
			},
		},
		{
			name:              "zero default expiration means no expiration",
			defaultExpiration: DefaultExpiration,
			expect: func(t *testing.T, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.True(expiration.IsZero())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(tc.defaultExpiration, 0)
			c.SetDefault(v1, v2)

			_, expiration, found := c.GetWithExpiration(v1)
			tc.expect(t, expiration, found)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(c Cache)
		expect func(t *testing.T, err error, x any, found bool)
	}{
		{
			name:  "add missing key",
			setup: func(c Cache) {},
			expect: func(t *testing.T, err error, x any, found bool) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(found)
				assert.Equal(v3, x)
			},
		},
		{
			name: "add existing key keeps the existing value",
			setup: func(c Cache) {
				c.Set(v1, v2, DefaultExpiration)
			},
			expect: func(t *testing.T, err error, x any, found bool) {
				assert := assert.New(t)
				assert.Error(err)
				assert.True(found)
				assert.Equal(v2, x)
			},
		},
		{
			name: "add expired key replaces it",
			setup: func(c Cache) {
				mockExpiredItem(c, v1, v2)
			},
			expect: func(t *testing.T, err error, x any, found bool) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(found)
				assert.Equal(v3, x)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(DefaultExpiration, 0)
			tc.setup(c)

			err := c.Add(v1, v3, DefaultExpiration)
			x, found := c.Get(v1)
			tc.expect(t, err, x, found)
		})
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	c.Set(v1, v2, DefaultExpiration)
	c.Delete(v1)

	x, found := c.Get(v1)
	assert.False(found)
	assert.Nil(x)
}

func TestCacheKeys(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	c.Set(v1, 1, DefaultExpiration)
	c.Set(v2, 2, DefaultExpiration)
	c.Set(v3, 3, DefaultExpiration)
	assert.ElementsMatch([]string{v1, v2, v3}, c.Keys())
}

func TestItems(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	mockExpiredItem(c, v1, "1")
	c.Set(v2, "2", DefaultExpiration)
	c.Set(v3, "3", DefaultExpiration)
	assert.Equal(map[string]Item{"bar": {Object: "2", Expiration: 0}, "baz": {Object: "3", Expiration: 0}}, c.Items())
}

func TestItemCount(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	c.Set(v1, "1", DefaultExpiration)
	c.Set(v2, "2", DefaultExpiration)
	c.Set(v3, "3", DefaultExpiration)
	assert.Equal(3, c.ItemCount())
}

func TestFlush(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	c.Set(v1, v2, DefaultExpiration)
	c.Set(v3, v4, DefaultExpiration)
	c.Flush()

	x, found := c.Get(v1)
	assert.False(found)
	assert.Nil(x)
	x, found = c.Get(v3)
	assert.False(found)
	assert.Nil(x)
}

func TestOnEvicted(t *testing.T) {
	assert := assert.New(t)
	c := New(DefaultExpiration, 0)
	c.Set(v1, 3, DefaultExpiration)

	var evictedKey string
	var evictedValue any
	c.OnEvicted(func(k string, v any) {
		evictedKey, evictedValue = k, v
		c.Set(v2, 4, DefaultExpiration)
	})
	c.Delete(v1)

	x, found := c.Get(v2)
	assert.Equal(v1, evictedKey)
	assert.Equal(3, evictedValue)
	assert.True(found)
	assert.Equal(4, x)
}

func TestGetWithExpiration(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		expect func(t *testing.T, x any, expiration time.Time, found bool)
	}{
		{
			name: "missing key",
			key:  "z",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.False(found)
				assert.Nil(x)
				assert.True(expiration.IsZero())
			},
		},
		{
			name: "int value with default expiration",
			key:  "a",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal(1, x)
				assert.True(expiration.IsZero())
			},
		},
		{
			name: "string value with default expiration",
			key:  "b",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal("b", x)
				assert.True(expiration.IsZero())
			},
		},
		{
			name: "float value with default expiration",
			key:  "c",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal(3.5, x)
				assert.True(expiration.IsZero())
			},
		},
		{
			name: "value with no expiration",
			key:  "d",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal(1, x)
				assert.True(expiration.IsZero())
			},
		},
		{
			name: "value with explicit expiration reports its deadline",
			key:  "e",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.True(found)
				assert.Equal(1, x)
				assert.True(expiration.After(time.Now()))
			},
		},
		{
			name: "expired value is not found",
			key:  "f",
			expect: func(t *testing.T, x any, expiration time.Time, found bool) {
				assert := assert.New(t)
				assert.False(found)
				assert.Nil(x)
				assert.True(expiration.IsZero())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(DefaultExpiration, 0)
			c.Set("a", 1, DefaultExpiration)
			c.Set("b", "b", DefaultExpiration)
			c.Set("c", 3.5, DefaultExpiration)
			c.Set("d", 1, NoExpiration)
			c.Set("e", 1, 50*time.Millisecond)
			mockExpiredItem(c, "f", 1)

			x, expiration, found := c.GetWithExpiration(tc.key)
			tc.expect(t, x, expiration, found)
		})
	}
}

func BenchmarkCacheGetExpiring(b *testing.B) {
	benchmarkCacheGet(b, 5*time.Minute)
}

func BenchmarkCacheGetNotExpiring(b *testing.B) {
	benchmarkCacheGet(b, NoExpiration)
}

func benchmarkCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(exp, 0)
	tc.Set(v1, v2, DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get(v1)
	}
}

func BenchmarkRWMutexMapGet(b *testing.B) {
	b.StopTimer()
	m := map[string]string{
		v1: v2,
	}
	var mu sync.RWMutex
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_, _ = m[v1]
		mu.RUnlock()
	}
}

func BenchmarkRWMutexInterfaceMapGetStruct(b *testing.B) {
	b.StopTimer()
	s := struct{ name string }{name: v1}
	m := map[any]string{
		s: v2,
	}
	var mu sync.RWMutex
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_, _ = m[s]
		mu.RUnlock()
	}
}

func BenchmarkRWMutexInterfaceMapGetString(b *testing.B) {
	b.StopTimer()
	m := map[any]string{
		v1: v2,
	}
	var mu sync.RWMutex
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_, _ = m[v1]
		mu.RUnlock()
	}
}

func BenchmarkCacheGetConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, NoExpiration)
}

func benchmarkCacheGetConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(exp, 0)
	tc.Set(v1, v2, DefaultExpiration)
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for range workers {
		go func() {
			for range each {
				tc.Get(v1)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkRWMutexMapGetConcurrent(b *testing.B) {
	b.StopTimer()
	m := map[string]string{
		v1: v2,
	}
	mu := sync.RWMutex{}
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for range workers {
		go func() {
			for range each {
				mu.RLock()
				_, _ = m[v1]
				mu.RUnlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, NoExpiration)
}

func benchmarkCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	n := 10000
	tc := New(exp, 0)
	keys := make([]string, n)
	for i := range n {
		k := v1 + strconv.Itoa(i)
		keys[i] = k
		tc.Set(k, v2, DefaultExpiration)
	}

	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for range each {
				tc.Get(k)
			}

			wg.Done()
		}(v)
	}

	b.StartTimer()
	wg.Wait()
}

func BenchmarkCacheSetExpiring(b *testing.B) {
	benchmarkCacheSet(b, 5*time.Minute)
}

func BenchmarkCacheSetNotExpiring(b *testing.B) {
	benchmarkCacheSet(b, NoExpiration)
}

func benchmarkCacheSet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(exp, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set(v1, v2, DefaultExpiration)
	}
}

func BenchmarkRWMutexMapSet(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[v1] = v2
		mu.Unlock()
	}
}

func BenchmarkCacheSetDelete(b *testing.B) {
	b.StopTimer()
	tc := New(DefaultExpiration, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set(v1, v2, DefaultExpiration)
		tc.Delete(v1)
	}
}

func BenchmarkRWMutexMapSetDelete(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[v1] = v2
		mu.Unlock()
		mu.Lock()
		delete(m, v1)
		mu.Unlock()
	}
}

func BenchmarkRWMutexMapSetDeleteSingleLock(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	var mu sync.Mutex
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[v1] = v2
		delete(m, v1)
		mu.Unlock()
	}
}
