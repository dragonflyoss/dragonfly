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
	"testing"
	"time"
)

func TestCacheConstants(t *testing.T) {
	// Test that cache size has been increased to handle high workloads
	if defaultCacheSize != 10240 {
		t.Errorf("Expected defaultCacheSize to be 10240, got %d", defaultCacheSize)
	}

	// Test that TTL has been increased to reduce cache churn
	expectedTTL := 5 * time.Minute
	if defaultTTL != expectedTTL {
		t.Errorf("Expected defaultTTL to be %v, got %v", expectedTTL, defaultTTL)
	}

	// Verify constants are reasonable for production use
	if defaultCacheSize < 1000 {
		t.Errorf("Cache size %d is too small for production workloads", defaultCacheSize)
	}

	if defaultTTL < time.Minute {
		t.Errorf("TTL %v is too short, may cause excessive database load", defaultTTL)
	}
}

func TestNewCacher(t *testing.T) {
	cacher, err := newCacher()
	if err != nil {
		t.Fatalf("Failed to create cacher: %v", err)
	}

	if cacher == nil {
		t.Fatal("Expected non-nil cacher")
	}

	// The cacher should implement the Cacher interface
	// We can't test internal fields but we can test it works
}
