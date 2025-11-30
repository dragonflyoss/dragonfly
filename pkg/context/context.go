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

package context

import (
	"context"
	"time"
)

const (
	// DefaultTimeout is the default timeout for background operations.
	DefaultTimeout = 30 * time.Second
)

// WithTimeout creates a new context with timeout from the parent context.
// If the parent context is nil or context.Background(), it creates a new context with timeout.
// This preserves tracing information and other context values from the parent.
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, timeout)
}

// WithDefaultTimeout creates a new context with default timeout from the parent context.
func WithDefaultTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(parent, DefaultTimeout)
}

// Copy creates a new context that copies values and tracing information from the parent,
// but with a new timeout. This is useful when you need to preserve context metadata
// but want independent timeout control.
func Copy(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		return context.WithTimeout(context.Background(), timeout)
	}

	// Create a new context with timeout, but preserve the parent's values
	// by using the parent as the base
	return context.WithTimeout(parent, timeout)
}

// CopyWithDefaultTimeout creates a new context with default timeout, copying from parent.
func CopyWithDefaultTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return Copy(parent, DefaultTimeout)
}
