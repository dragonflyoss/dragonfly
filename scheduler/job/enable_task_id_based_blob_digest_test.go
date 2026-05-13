/*
 *     Copyright 2025 The Dragonfly Authors
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

package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnableTaskIDBasedBlobDigest(t *testing.T) {
	tests := []struct {
		name   string
		enable *bool
		expect bool
	}{
		{
			name:   "nil returns false (default for backward compatibility)",
			enable: nil,
			expect: false,
		},
		{
			name:   "true returns true",
			enable: func() *bool { b := true; return &b }(),
			expect: true,
		},
		{
			name:   "false returns false",
			enable: func() *bool { b := false; return &b }(),
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getEnableTaskIDBasedBlobDigest(tc.enable)
			assert.Equal(t, tc.expect, result)
		})
	}
}
