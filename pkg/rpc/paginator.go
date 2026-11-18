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

var (
	// ErrInvalidPageToken represents the page token is invalid.
	ErrInvalidPageToken = errors.New("invalid page token")
)

// Tokenizer returns the key of the item. The paginator orders the items in
// the group by the key and encodes the key of the last item in the page to the
// page token, so the key must be unique in the group and not empty.
type Tokenizer[T any] func(item T) string

// Paginator pages the items in the order of the groups, and in the order of
// the keys in the group, refer to
// https://github.com/hashicorp/nomad/blob/main/nomad/state/paginator/paginator.go.
type Paginator[T any] struct {
	// pageSize is the maximum number of the items in the page, the page has
	// all the remaining items if it is zero.
	pageSize int32

	// group is the index of the group of the last item in the previous page.
	group int

	// start is the key of the last item in the previous page, the page starts
	// after it.
	start string

	// tokenizer returns the key of the item.
	tokenizer Tokenizer[T]
}

// pageToken is the content of the page token. The page token is opaque to the
// caller, so the content can be changed with a new version.
type pageToken struct {
	// Version is the version of the page token.
	Version string `json:"v"`

	// Group is the index of the group of the last item in the previous page.
	Group int `json:"group"`

	// Start is the key of the last item in the previous page, the next page
	// starts after it.
	Start string `json:"start"`
}

// keyedItem pairs an item with its group and pre-computed key.
type keyedItem[T any] struct {
	// group is the index of the group of the item.
	group int

	// key is the key used to order the item in the group.
	key string

	// item is the item.
	item T
}

// NewPaginator returns a Paginator instance. The error is caused by the
// invalid page token of the caller, so the gRPC handler should return the
// invalid argument error.
func NewPaginator[T any](pageSize int32, pageToken string, tokenizer Tokenizer[T]) (*Paginator[T], error) {
	p := &Paginator[T]{
		pageSize:  pageSize,
		tokenizer: tokenizer,
	}

	// Start from the first item if the page token is not specified.
	if pageToken == "" {
		return p, nil
	}

	group, start, err := p.decodePageToken(pageToken)
	if err != nil {
		return nil, err
	}

	p.group = group
	p.start = start
	return p, nil
}

// Page returns the page of the items, and the token of the next page if there
// are more items. The items of a group are before the items of the following
// groups, and the items in the group are in the order of the keys.
func (p *Paginator[T]) Page(groups ...[]T) ([]T, *string, error) {
	var count int
	for _, items := range groups {
		count += len(items)
	}

	// Calculate the key once for each item, and skip the items in the previous pages.
	keyedItems := make([]keyedItem[T], 0, count)
	for group, items := range groups {
		for _, item := range items {
			key := p.tokenizer(item)
			if cmp.Or(cmp.Compare(group, p.group), cmp.Compare(key, p.start)) > 0 {
				keyedItems = append(keyedItems, keyedItem[T]{group: group, key: key, item: item})
			}
		}
	}

	slices.SortFunc(keyedItems, func(a, b keyedItem[T]) int {
		return cmp.Or(cmp.Compare(a.group, b.group), cmp.Compare(a.key, b.key))
	})

	// Truncate the items to the page size and set the next page token if
	// there are more items than the page size.
	var nextPageToken *string
	if p.pageSize > 0 && len(keyedItems) > int(p.pageSize) {
		keyedItems = keyedItems[:p.pageSize]
		last := keyedItems[p.pageSize-1]
		token, err := p.encodePageToken(last.group, last.key)
		if err != nil {
			return nil, nil, err
		}

		nextPageToken = &token
	}

	page := make([]T, 0, len(keyedItems))
	for _, keyed := range keyedItems {
		page = append(page, keyed.item)
	}

	return page, nextPageToken, nil
}

// encodePageToken encodes the group and the key of the last item in the page
// to the page token, refer to https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apiserver/pkg/storage/continue.go.
func (p *Paginator[T]) encodePageToken(group int, start string) (string, error) {
	data, err := json.Marshal(&pageToken{Version: pageTokenVersion, Group: group, Start: start})
	if err != nil {
		return "", fmt.Errorf("marshal page token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

// decodePageToken decodes the page token to the group and the key of the last
// item in the previous page.
func (p *Paginator[T]) decodePageToken(token string) (int, string, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, "", fmt.Errorf("%w: %w", ErrInvalidPageToken, err)
	}

	var pt pageToken
	if err := json.Unmarshal(data, &pt); err != nil {
		return 0, "", fmt.Errorf("%w: %w", ErrInvalidPageToken, err)
	}

	if pt.Version != pageTokenVersion {
		return 0, "", fmt.Errorf("%w: unsupported version %q", ErrInvalidPageToken, pt.Version)
	}

	if pt.Group < 0 {
		return 0, "", fmt.Errorf("%w: negative group %d", ErrInvalidPageToken, pt.Group)
	}

	if pt.Start == "" {
		return 0, "", fmt.Errorf("%w: empty start key", ErrInvalidPageToken)
	}

	return pt.Group, pt.Start, nil
}
