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

package balancer

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

type mockSubConn struct {
	balancer.SubConn
	id string
}

func mockReadySCs(subConns ...*mockSubConn) map[balancer.SubConn]base.SubConnInfo {
	readySCs := make(map[balancer.SubConn]base.SubConnInfo, len(subConns))
	for _, sc := range subConns {
		readySCs[sc] = base.SubConnInfo{Address: resolver.Address{Addr: sc.id + ":8002", ServerName: sc.id}}
	}

	return readySCs
}

func TestConsistentHashingPickerBuilder_Build(t *testing.T) {
	foo, bar, baz := &mockSubConn{id: "foo"}, &mockSubConn{id: "bar"}, &mockSubConn{id: "baz"}

	tests := []struct {
		name     string
		readySCs map[balancer.SubConn]base.SubConnInfo
		expect   func(t *testing.T, picker balancer.Picker)
	}{
		{
			name: "no ready subconns returns error picker",
			expect: func(t *testing.T, picker balancer.Picker) {
				assert := assert.New(t)
				result, err := picker.Pick(balancer.PickInfo{Ctx: context.WithValue(context.Background(), ContextKey, "task")})
				assert.ErrorIs(err, balancer.ErrNoSubConnAvailable)
				assert.Nil(result.SubConn)
			},
		},
		{
			name:     "ready subconns are keyed by address and server name",
			readySCs: mockReadySCs(foo, bar, baz),
			expect: func(t *testing.T, picker balancer.Picker) {
				assert := assert.New(t)
				p, ok := picker.(*consistentHashingPicker)
				if !assert.True(ok) {
					return
				}

				assert.Equal(map[string]balancer.SubConn{
					"foo:8002:foo": foo,
					"bar:8002:bar": bar,
					"baz:8002:baz": baz,
				}, p.subConns)
				assert.ElementsMatch([]string{"foo:8002:foo", "bar:8002:bar", "baz:8002:baz"}, p.hashring.Members())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &ConsistentHashingPickerBuilder{}
			tc.expect(t, b.Build(base.PickerBuildInfo{ReadySCs: tc.readySCs}))
		})
	}
}

func TestConsistentHashingPicker_Pick(t *testing.T) {
	b := &ConsistentHashingPickerBuilder{}
	picker := b.Build(base.PickerBuildInfo{ReadySCs: mockReadySCs(&mockSubConn{id: "foo"}, &mockSubConn{id: "bar"}, &mockSubConn{id: "baz"})}).(*consistentHashingPicker)

	tests := []struct {
		name   string
		ctx    context.Context
		expect func(t *testing.T, result balancer.PickResult, err error)
	}{
		{
			name: "context without task id",
			ctx:  context.Background(),
			expect: func(t *testing.T, result balancer.PickResult, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(result.SubConn)
			},
		},
		{
			name: "task id is not a string",
			ctx:  context.WithValue(context.Background(), ContextKey, 1),
			expect: func(t *testing.T, result balancer.PickResult, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(result.SubConn)
			},
		},
		{
			name: "task foo picks hashring member",
			ctx:  context.WithValue(context.Background(), ContextKey, "foo"),
			expect: func(t *testing.T, result balancer.PickResult, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				element, err := picker.hashring.Get("foo")
				assert.NoError(err)
				assert.Same(picker.subConns[element], result.SubConn)
			},
		},
		{
			name: "task bar picks hashring member",
			ctx:  context.WithValue(context.Background(), ContextKey, "bar"),
			expect: func(t *testing.T, result balancer.PickResult, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				element, err := picker.hashring.Get("bar")
				assert.NoError(err)
				assert.Same(picker.subConns[element], result.SubConn)
			},
		},
		{
			name: "task baz picks hashring member",
			ctx:  context.WithValue(context.Background(), ContextKey, "baz"),
			expect: func(t *testing.T, result balancer.PickResult, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				element, err := picker.hashring.Get("baz")
				assert.NoError(err)
				assert.Same(picker.subConns[element], result.SubConn)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := picker.Pick(balancer.PickInfo{Ctx: tc.ctx})
			tc.expect(t, result, err)
		})
	}
}

func TestConsistentHashingPicker_Pick_Deterministic(t *testing.T) {
	assert := assert.New(t)
	readySCs := mockReadySCs(&mockSubConn{id: "foo"}, &mockSubConn{id: "bar"}, &mockSubConn{id: "baz"})
	picker := (&ConsistentHashingPickerBuilder{}).Build(base.PickerBuildInfo{ReadySCs: readySCs})
	rebuilt := (&ConsistentHashingPickerBuilder{}).Build(base.PickerBuildInfo{ReadySCs: readySCs})

	for i := range 20 {
		ctx := context.WithValue(context.Background(), ContextKey, fmt.Sprintf("task-%d", i))
		first, err := picker.Pick(balancer.PickInfo{Ctx: ctx})
		assert.NoError(err)
		second, err := picker.Pick(balancer.PickInfo{Ctx: ctx})
		assert.NoError(err)
		assert.Same(first.SubConn, second.SubConn)

		other, err := rebuilt.Pick(balancer.PickInfo{Ctx: ctx})
		assert.NoError(err)
		assert.Same(first.SubConn, other.SubConn)
	}
}

func TestConsistentHashingPickerBuilder_GetCircle(t *testing.T) {
	tests := []struct {
		name     string
		readySCs map[balancer.SubConn]base.SubConnInfo
		expect   func(t *testing.T, b *ConsistentHashingPickerBuilder, circle map[string]string, err error)
	}{
		{
			name: "before build",
			expect: func(t *testing.T, b *ConsistentHashingPickerBuilder, circle map[string]string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(circle)
			},
		},
		{
			name:     "single member",
			readySCs: mockReadySCs(&mockSubConn{id: "foo"}),
			expect: func(t *testing.T, b *ConsistentHashingPickerBuilder, circle map[string]string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.ElementsMatch(b.hashring.Members(), slices.Collect(maps.Keys(circle)))
				for member, key := range circle {
					got, err := b.hashring.Get(key)
					assert.NoError(err)
					assert.Equal(member, got)
				}
			},
		},
		{
			name:     "multiple members",
			readySCs: mockReadySCs(&mockSubConn{id: "foo"}, &mockSubConn{id: "bar"}, &mockSubConn{id: "baz"}),
			expect: func(t *testing.T, b *ConsistentHashingPickerBuilder, circle map[string]string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.ElementsMatch(b.hashring.Members(), slices.Collect(maps.Keys(circle)))
				for member, key := range circle {
					got, err := b.hashring.Get(key)
					assert.NoError(err)
					assert.Equal(member, got)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &ConsistentHashingPickerBuilder{}
			if tc.readySCs != nil {
				b.Build(base.PickerBuildInfo{ReadySCs: tc.readySCs})
			}

			circle, err := b.GetCircle()
			tc.expect(t, b, circle, err)
		})
	}
}

func TestConsistentHashingPickerBuilder_GetCircle_Cached(t *testing.T) {
	assert := assert.New(t)
	b := &ConsistentHashingPickerBuilder{}
	b.Build(base.PickerBuildInfo{ReadySCs: mockReadySCs(&mockSubConn{id: "foo"})})

	first, err := b.GetCircle()
	assert.NoError(err)
	second, err := b.GetCircle()
	assert.NoError(err)
	assert.Equal(reflect.ValueOf(first).Pointer(), reflect.ValueOf(second).Pointer())

	b.Build(base.PickerBuildInfo{ReadySCs: mockReadySCs(&mockSubConn{id: "foo"}, &mockSubConn{id: "bar"})})
	rebuilt, err := b.GetCircle()
	assert.NoError(err)
	assert.NotEqual(reflect.ValueOf(first).Pointer(), reflect.ValueOf(rebuilt).Pointer())
	assert.Len(rebuilt, 2)
}
