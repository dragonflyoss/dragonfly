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
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePersonalAccessToken(t *testing.T) {
	assert := assert.New(t)

	token := GeneratePersonalAccessToken()
	assert.NotEmpty(token)

	// Token must be url-safe base64 wrapping a valid UUID.
	raw, err := base64.RawURLEncoding.DecodeString(token)
	assert.NoError(err)
	_, err = uuid.Parse(string(raw))
	assert.NoError(err)

	// Tokens must be unique.
	assert.NotEqual(token, GeneratePersonalAccessToken())
}
