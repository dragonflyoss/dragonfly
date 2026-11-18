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

package rpc

import (
	"cmp"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

const (
	// pageTokenVersion is the version of the page token.
	pageTokenVersion = "v1"
)

// pageToken is the content of the page token. The page token is opaque to the
// caller, so the content can be changed with a new version.
type pageToken struct {
	// Version is the version of the page token.
	Version string `json:"v"`

	// Start is the key of the last item in the previous page, the next page
	// starts after it.
	Start string `json:"start"`
}

// keyedItem is the item with the key used to order the items.
type keyedItem[T any] struct {
	// key is the key of the item.
	key string

	// item is the item.
	item T
}

// encodePageToken encodes the key of the last item in the page to the page
// token, refer to https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apiserver/pkg/storage/continue.go.
func encodePageToken(start string) (string, error) {
	data, err := json.Marshal(&pageToken{Version: pageTokenVersion, Start: start})
	if err != nil {
		return "", fmt.Errorf("marshal page token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

// decodePageToken decodes the page token to the key of the last item in the
// previous page.
func decodePageToken(token string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("decode page token: %w", err)
	}

	var pt pageToken
	if err := json.Unmarshal(data, &pt); err != nil {
		return "", fmt.Errorf("unmarshal page token: %w", err)
	}

	if pt.Version != pageTokenVersion {
		return "", fmt.Errorf("unsupported page token version %q", pt.Version)
	}

	if pt.Start == "" {
		return "", errors.New("page token has no start key")
	}

	return pt.Start, nil
}

// Paginate returns the page of the items in the order of the keys, and the
// token of the next page if there are more items. The page starts after the
// key in the page token, and has at most pageSize items. If pageSize is zero,
// all the remaining items are returned.
func Paginate[T any](items []T, key func(T) string, pageSize int32, token string) ([]T, *string, error) {
	var start string
	if token != "" {
		var err error
		if start, err = decodePageToken(token); err != nil {
			return nil, nil, err
		}
	}

	// Calculate the key once for each item, and skip the items in the previous pages.
	keyedItems := make([]keyedItem[T], 0, len(items))
	for _, item := range items {
		if k := key(item); k > start {
			keyedItems = append(keyedItems, keyedItem[T]{key: k, item: item})
		}
	}

	slices.SortFunc(keyedItems, func(a, b keyedItem[T]) int {
		return cmp.Compare(a.key, b.key)
	})

	// Set the next page token if there are more items than the page size.
	var nextPageToken *string
	if pageSize > 0 && len(keyedItems) > int(pageSize) {
		keyedItems = keyedItems[:pageSize]
		token, err := encodePageToken(keyedItems[len(keyedItems)-1].key)
		if err != nil {
			return nil, nil, err
		}

		nextPageToken = &token
	}

	page := make([]T, 0, len(keyedItems))
	for _, keyedItem := range keyedItems {
		page = append(page, keyedItem.item)
	}

	return page, nextPageToken, nil
}
