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

package ping

import (
	"errors"
	"time"

	"github.com/gaius-qi/ping"
)

// Ping sends
func Ping(addr string) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return nil, err
	}

	pinger.Count = 1
	pinger.Timeout = 100 * time.Millisecond
	err = pinger.Run()
	if err != nil {
		return nil, err
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		return nil, errors.New("receive packet failed")
	}

	return stats, nil
}
