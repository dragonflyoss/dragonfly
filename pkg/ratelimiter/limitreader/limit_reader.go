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

package limitreader

import (
	"fmt"
	"hash"
	"io"

	"d7y.io/dragonfly/v2/internal/constants"
	"github.com/opencontainers/go-digest"

	"d7y.io/dragonfly/v2/pkg/ratelimiter/ratelimiter"
	"d7y.io/dragonfly/v2/pkg/util/digestutils"
)

// NewLimitReader creates a LimitReader.
// src: reader
// rate: bytes/second
func NewLimitReader(src io.Reader, rate int64) *LimitReader {
	return NewLimitReaderWithLimiter(newRateLimiterWithDefaultWindow(rate), src)
}

// NewLimitReaderWithLimiter creates LimitReader with a rateLimiter.
// src: reader
// rate: bytes/second
func NewLimitReaderWithLimiter(rl *ratelimiter.RateLimiter, src io.Reader) *LimitReader {
	return &LimitReader{
		Src:     src,
		Limiter: rl,
	}
}

// NewLimitReaderWithMD5Sum creates LimitReader with a md5 sum.
// src: reader
// rate: bytes/second
func NewLimitReaderWithMD5Sum(src io.Reader, rate int64, md5sum hash.Hash) *LimitReader {
	return NewLimitReaderWithLimiterAndMD5Sum(src, newRateLimiterWithDefaultWindow(rate), md5sum)
}

// NewLimitReaderWithLimiterAndMD5Sum creates LimitReader with rateLimiter and md5 sum.
// src: reader
// rate: bytes/second
func NewLimitReaderWithLimiterAndMD5Sum(src io.Reader, rl *ratelimiter.RateLimiter, md5sum hash.Hash) *LimitReader {
	return &LimitReader{
		Src:        src,
		Limiter:    rl,
		digest:     md5sum,
		digestType: constants.Md5Hash.String(),
	}
}

func newRateLimiterWithDefaultWindow(rate int64) *ratelimiter.RateLimiter {
	return ratelimiter.NewRateLimiter(ratelimiter.TransRate(rate), 2)
}

// LimitReader reads stream with RateLimiter.
type LimitReader struct {
	Src        io.Reader
	Limiter    *ratelimiter.RateLimiter
	digest     hash.Hash
	digestType string
}

func (lr *LimitReader) Read(p []byte) (n int, err error) {
	n, e := lr.Src.Read(p)
	if e != nil && e != io.EOF {
		return n, e
	}
	if n > 0 {
		if lr.digest != nil {
			lr.digest.Write(p[:n])
		}
		lr.Limiter.AcquireBlocking(int64(n))
	}
	return n, e
}

// Md5 calculates the md5 of all contents read.
func (lr *LimitReader) Md5() string {
	if lr.digest != nil {
		return digestutils.ToHashString(lr.digest)
	}
	return ""
}

// NewLimitReaderWithLimiterAndDigest creates LimitReader with rateLimiter and digest.
// src: reader
// rate: bytes/second
func NewLimitReaderWithLimiterAndDigest(src io.Reader, rl *ratelimiter.RateLimiter, digest hash.Hash, digestType digest.Algorithm) *LimitReader {
	return &LimitReader{
		Src:        src,
		Limiter:    rl,
		digest:     digest,
		digestType: digestType.String(),
	}
}

// Digest calculates the digest of all contents read, return value is like <algo>:<hex_value>
func (lr *LimitReader) Digest() string {
	if lr.digest != nil {
		return fmt.Sprintf("%s:%s", lr.digestType, digestutils.ToHashString(lr.digest))
	}
	return ""
}
