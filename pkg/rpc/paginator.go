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
	"container/heap"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
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

	// Keep one more item than the page size to know if there are more items,
	// and keep all the items if the page size is not specified.
	limit := count
	if p.pageSize > 0 {
		limit = min(count, int(p.pageSize)+1)
	}

	start := keyedItem[T]{group: p.group, key: p.start}
	keyedItems := &keyedItemHeap[T]{items: make([]keyedItem[T], 0, limit)}
	for group, items := range groups {
		for _, item := range items {
			// Calculate the key once for each item, and skip the items in the previous pages.
			keyed := keyedItem[T]{group: group, key: p.tokenizer(item), item: item}
			if keyed.compare(&start) <= 0 {
				continue
			}

			// Build a heap of up to limit items with the last item at the top, so the
			// page costs the memory of the page size rather than the count of the items, refer to
			// https://github.com/prometheus/prometheus/blob/main/promql/engine.go.
			switch {
			case p.pageSize == 0:
				keyedItems.items = append(keyedItems.items, keyed)
			case len(keyedItems.items) < limit:
				heap.Push(keyedItems, keyed)
			case keyed.compare(&keyedItems.items[0]) < 0:
				keyedItems.items[0] = keyed
				heap.Fix(keyedItems, 0)
			}
		}
	}

	slices.SortFunc(keyedItems.items, func(a, b keyedItem[T]) int {
		return a.compare(&b)
	})

	// Truncate the items to the page size and set the next page token if
	// there are more items than the page size.
	var nextPageToken *string
	if p.pageSize > 0 && len(keyedItems.items) > int(p.pageSize) {
		keyedItems.items = keyedItems.items[:p.pageSize]
		last := keyedItems.items[p.pageSize-1]
		token, err := p.encodePageToken(last.group, last.key)
		if err != nil {
			return nil, nil, err
		}

		nextPageToken = &token
	}

	page := make([]T, 0, len(keyedItems.items))
	for _, keyed := range keyedItems.items {
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

// keyedItemHeap is the max-heap of the keyed items, the top of the heap is the
// last item in the order of the page, refer to
// https://github.com/pingcap/tidb/blob/master/pkg/executor/sortexec/topn_chunk_heap.go.
type keyedItemHeap[T any] struct {
	// items are the keyed items in the heap.
	items []keyedItem[T]
}

// compare compares the item with the other item in the order of the page, the
// group is compared before the key.
func (k *keyedItem[T]) compare(other *keyedItem[T]) int {
	if order := cmp.Compare(k.group, other.group); order != 0 {
		return order
	}

	return strings.Compare(k.key, other.key)
}

// Len returns the length of the heap.
func (h *keyedItemHeap[T]) Len() int {
	return len(h.items)
}

// Less reports whether the item with index i is after the item with index j
// in the order of the page.
func (h *keyedItemHeap[T]) Less(i, j int) bool {
	return h.items[i].compare(&h.items[j]) > 0
}

// Swap swaps the items with indexes i and j.
func (h *keyedItemHeap[T]) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Push pushes the item to the heap.
func (h *keyedItemHeap[T]) Push(item any) {
	h.items = append(h.items, item.(keyedItem[T]))
}

// Pop pops the last item from the heap.
func (h *keyedItemHeap[T]) Pop() any {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}
