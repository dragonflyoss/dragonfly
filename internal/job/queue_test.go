/*
 *     Copyright 2020 The Dragonfly Authors
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

func TestJob_GetSchedulerQueue(t *testing.T) {
	tests := []struct {
		name      string
		clusterID uint
		hostname  string
		ip        string
		expect    func(t *testing.T, result Queue, err error)
	}{
		{
			name:      "get scheduler queue",
			clusterID: 1,
			hostname:  "foo",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, result Queue, err error) {
				assert := assert.New(t)
				assert.Equal(Queue("scheduler_1_foo_127.0.0.1"), result)
			},
		},
		{
			name:      "get scheduler queue with empty hostname",
			clusterID: 1,
			hostname:  "",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, result Queue, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "empty hostname config is not specified")
			},
		},
		{
			name:      "get scheduler queue with empty clusterID",
			clusterID: 0,
			hostname:  "foo",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, result Queue, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "empty cluster id config is not specified")
			},
		},
		{
			name:      "get scheduler queue with empty ip",
			clusterID: 1,
			hostname:  "foo",
			ip:        "",
			expect: func(t *testing.T, result Queue, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "empty ip config is not specified")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			queue, err := GetSchedulerQueue(tc.clusterID, tc.hostname, tc.ip)
			tc.expect(t, queue, err)
		})
	}
}
