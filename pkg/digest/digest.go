/*
 *     Copyright 2022 The Dragonfly Authors
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

package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// AlgorithmCRC32 is crc32 algorithm name of hash.
	AlgorithmCRC32 = "crc32"

	// AlgorithmBlake3 is blake3 algorithm name of hash.
	AlgorithmBlake3 = "blake3"

	// AlgorithmSHA1 is sha1 algorithm name of hash.
	AlgorithmSHA1 = "sha1"

	// AlgorithmSHA256 is sha256 algorithm name of hash.
	AlgorithmSHA256 = "sha256"

	// AlgorithmSHA512 is sha512 algorithm name of hash.
	AlgorithmSHA512 = "sha512"

	// AlgorithmMD5 is md5 algorithm name of hash.
	AlgorithmMD5 = "md5"
)

// Digest provides digest operation function.
type Digest struct {
	// Algorithm is hash algorithm.
	Algorithm string

	// Encoded is hash encode.
	Encoded string
}

// String return digest string.
func (d *Digest) String() string {
	return fmt.Sprintf("%s:%s", d.Algorithm, d.Encoded)
}

// New return digest instance.
func New(algorithm, encoded string) *Digest {
	return &Digest{
		Algorithm: algorithm,
		Encoded:   encoded,
	}
}

// Parse uses to parse digest string to algorithm and encoded.
func Parse(digest string) (*Digest, error) {
	values := strings.Split(strings.TrimSpace(digest), ":")
	if len(values) != 2 {
		return nil, errors.New("invalid digest")
	}

	algorithm := values[0]
	encoded := values[1]

	switch algorithm {
	case AlgorithmCRC32:
		if len(encoded) <= 0 {
			return nil, errors.New("invalid encoded")
		}
	case AlgorithmBlake3:
		if len(encoded) != 64 {
			return nil, errors.New("invalid encoded")
		}
	case AlgorithmSHA1:
		if len(encoded) != 40 {
			return nil, errors.New("invalid encoded")
		}
	case AlgorithmSHA256:
		if len(encoded) != 64 {
			return nil, errors.New("invalid encoded")
		}
	case AlgorithmSHA512:
		if len(encoded) != 128 {
			return nil, errors.New("invalid encoded")
		}
	case AlgorithmMD5:
		if len(encoded) != 32 {
			return nil, errors.New("invalid encoded")
		}
	default:
		return nil, errors.New("invalid algorithm")
	}

	// The supported algorithms all produce hex-encoded digests. The switch only
	// checks the length, so reject any encoded value that is not valid
	// hexadecimal, otherwise a malformed digest such as a 64 character non-hex
	// string would be accepted as a sha256.
	if !isHex(encoded) {
		return nil, errors.New("invalid encoded")
	}

	return &Digest{
		Algorithm: algorithm,
		Encoded:   encoded,
	}, nil
}

// isHex reports whether s consists solely of hexadecimal characters.
func isHex(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}

	return true
}

// SHA256FromStrings computes the SHA256 checksum with multiple strings.
func SHA256FromStrings(data ...string) string {
	if len(data) == 0 {
		return ""
	}

	h := sha256.New()
	for _, s := range data {
		if _, err := h.Write([]byte(s)); err != nil {
			return ""
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

// blobURLRegex is the regex pattern for OCI blob URLs,
// e.g. http(s)://<registry>/v2/<repository>/blobs/<digest>.
var blobURLRegex = regexp.MustCompile(`^(.*)://(.*)/v2/(.*)/blobs/([^?]+)(?:\?.*)?$`)

// IsBlobURL checks whether the url is an OCI blob URL.
func IsBlobURL(url string) bool {
	return blobURLRegex.MatchString(url)
}

// ExtractFromBlobURL extracts the digest from an OCI blob URL.
// e.g. https://registry.example.com/v2/library/nginx/blobs/sha256:abc123...
//
//	returns Digest{Algorithm: "sha256", Encoded: "abc123..."}, nil
func ExtractFromBlobURL(url string) (*Digest, error) {
	matches := blobURLRegex.FindStringSubmatch(url)
	if len(matches) < 5 {
		return nil, errors.New("invalid blob URL")
	}

	return Parse(matches[4])
}
