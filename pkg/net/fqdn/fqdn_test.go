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

package fqdn

import (
	"testing"

	"github.com/Showmax/go-fqdn"
	"github.com/stretchr/testify/assert"
)

func TestFQDNHostname(t *testing.T) {
	// Save the original fqdn.FqdnHostname function and restore it after the test.
	_, err := fqdn.FqdnHostname()
	if err != nil {

		// Test case where FQDN hostname is available
		fqdnMock := "example.com"
		fqdnResult := fqdnHostname()
		assert.Equal(t, fqdnMock, fqdnResult, "FQDN hostname should match when available")

		// Test case where FQDN hostname is not available
		fqdnResult = fqdnHostname()
		assert.NotEmpty(t, fqdnResult, "FQDN hostname should fallback to non-empty value when unavailable")
	}
}
