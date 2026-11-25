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
	"time"
)

// defaultRefreshRatio is the fraction of a token's lifetime after which the
// cached token is refreshed. At 0.8, a 10-minute token is re-issued on the
// first call after the 8-minute mark, keeping the last 20% as a safety margin
// against clock skew between components.
const defaultRefreshRatio = 0.8

// TokenProvider issues tokens for outgoing calls and caches them until they
// approach expiry, so that high-frequency callers (dynconfig polling, announcer
// heartbeats) do not sign a new token on every RPC. Construct it with
// NewTokenProvider.
type TokenProvider struct {
	tm *TokenManager

	// subject is the subject of the issued tokens.
	subject string

	// ttl is the lifetime of the issued tokens.
	ttl time.Duration

	// refreshAt is the token age at which the cached token is re-issued.
	refreshAt time.Duration

	// mu guards cached. It is held across the signing call so that
	// concurrent callers block and share the result of a single refresh
	// instead of each signing their own token.
	mu sync.Mutex

	// cached is the token currently held; nil before the first issue.
	cached *cachedToken

	// now is the clock, overridable for tests.
	now func() time.Time
}

// cachedToken is the token currently held by a TokenProvider.
type cachedToken struct {
	value     string
	issuedAt  time.Time
	expiresAt time.Time
}

// NewTokenProvider returns a TokenProvider that issues tokens for subject with
// the given ttl and re-issues the cached token once defaultRefreshRatio of the
// lifetime has elapsed. A non-positive ttl falls back to DefaultTokenTTL.
func NewTokenProvider(tm *TokenManager, subject string, ttl time.Duration) *TokenProvider {
	if ttl <= 0 {
		ttl = DefaultTokenTTL
	}

	return &TokenProvider{
		tm:        tm,
		subject:   subject,
		ttl:       ttl,
		refreshAt: time.Duration(float64(ttl) * defaultRefreshRatio),
		now:       time.Now,
	}
}

// Token returns a token for an outgoing call, issuing a new one when the
// cached token is missing or has lived past its refresh point. It is safe for
// concurrent use: a refresh is serialized under the provider's lock, so
// concurrent callers share one signature instead of each producing their own.
func (p *TokenProvider) Token() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cached != nil {
		age := p.now().Sub(p.cached.issuedAt)
		if age < p.refreshAt && age < p.ttl {
			return p.cached.value, nil
		}
	}

	value, err := p.tm.IssueToken(p.subject, p.ttl)
	if err != nil {
		return "", err
	}

	p.cached = &cachedToken{
		value:     value,
		issuedAt:  p.now(),
		expiresAt: p.now().Add(p.ttl),
	}
	return value, nil
}
