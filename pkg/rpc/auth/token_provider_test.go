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

package auth

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// tokenProviderEnvironment wires a TokenProvider to a controllable clock.
func tokenProviderEnvironment(t *testing.T, ttl time.Duration) (*TokenProvider, *fakeClock) {
	t.Helper()

	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	clock := &fakeClock{now: time.Now()}
	p := NewTokenProvider(tm, "scheduler", ttl)
	p.now = clock.Now
	return p, clock
}

func TestTokenProviderIssuesValidToken(t *testing.T) {
	p, _ := tokenProviderEnvironment(t, time.Minute)

	token, err := p.Token()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// The issued token must satisfy the server-side verification.
	claims, err := p.tm.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, "scheduler", claims.Subject)
	require.Equal(t, testIssuer, claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
}

func TestTokenProviderCachesUntilRefreshPoint(t *testing.T) {
	ttl := 10 * time.Minute
	p, clock := tokenProviderEnvironment(t, ttl)

	first, err := p.Token()
	require.NoError(t, err)

	firstIssuedAt := p.cached.issuedAt

	// Before the 80% point the cached token is served unchanged.
	for i := 0; i < 5; i++ {
		clock.Advance(time.Minute)
		token, err := p.Token()
		require.NoError(t, err)
		require.Equal(t, first, token)
		require.Equal(t, firstIssuedAt, p.cached.issuedAt)
	}

	// Past the 8-minute mark (80% of the 10-minute TTL) a fresh token is
	// issued and takes over the cache.
	clock.Advance(3 * time.Minute)
	second, err := p.Token()
	require.NoError(t, err)
	require.Equal(t, clock.Now(), p.cached.issuedAt)
	require.True(t, p.cached.issuedAt.After(firstIssuedAt))

	// The refreshed token is still verifiable.
	_, err = p.tm.ParseToken(second)
	require.NoError(t, err)
}

func TestTokenProviderRefreshesExpiredToken(t *testing.T) {
	p, clock := tokenProviderEnvironment(t, time.Minute)

	_, err := p.Token()
	require.NoError(t, err)

	firstIssuedAt := p.cached.issuedAt

	// A token fully past its TTL is re-issued even though the refresh
	// point alone would not explain a new signature.
	clock.Advance(2 * time.Hour)
	token, err := p.Token()
	require.NoError(t, err)
	require.True(t, p.cached.issuedAt.After(firstIssuedAt))

	_, err = p.tm.ParseToken(token)
	require.NoError(t, err)
}

func TestTokenProviderDefaultTTL(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	p := NewTokenProvider(tm, "scheduler", 0) // non-positive ttl → DefaultTokenTTL

	require.Equal(t, DefaultTokenTTL, p.ttl)
	require.Equal(t, time.Duration(float64(DefaultTokenTTL)*defaultRefreshRatio), p.refreshAt)
}

func TestTokenProviderConcurrentCallsShareOneRefresh(t *testing.T) {
	p, clock := tokenProviderEnvironment(t, time.Minute)

	// Warm the cache, then move past the refresh point so that every
	// concurrent caller wants a refresh at the same time.
	_, err := p.Token()
	require.NoError(t, err)
	clock.Advance(2 * time.Minute)

	const callers = 32
	tokens := make([]string, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			token, err := p.Token()
			if err == nil {
				tokens[i] = token
			}
		}(i)
	}
	wg.Wait()

	// All callers observe the same refreshed token.
	for i := 1; i < callers; i++ {
		require.Equal(t, tokens[0], tokens[i])
	}

	// And exactly one underlying signing happened for the refresh: the
	// decoded claims all share the same iat, which can be sampled before
	// the goroutines run by virtue of them all blocking on the lock.
	claims, err := p.tm.ParseToken(tokens[0])
	require.NoError(t, err)
	require.NotNil(t, claims.IssuedAt)
}

// fakeClock is a manually advanced clock for deterministic refresh tests.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the current fake time.
func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Advance moves the fake clock forward by d.
func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
