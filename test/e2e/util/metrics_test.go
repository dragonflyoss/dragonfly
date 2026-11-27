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

package util

import "testing"

func TestParseDownloadTraffic(t *testing.T) {
	const counter = "# TYPE dragonfly_client_download_traffic counter\n"
	tests := []struct {
		name    string
		input   string
		want    DownloadTraffic
		wantErr bool
	}{
		{
			name: "separate traffic types",
			input: counter + `dragonfly_client_download_traffic{type="REMOTE_PEER"} 33554432
dragonfly_client_download_traffic{type="BACK_TO_SOURCE"} 4096
dragonfly_client_download_traffic{type="LOCAL_PEER"} 8192
`,
			want: DownloadTraffic{RemotePeer: 33554432, BackToSource: 4096},
		},
		{
			name: "aggregate label sets by type",
			input: counter + `dragonfly_client_download_traffic{type="REMOTE_PEER",application="one"} 1024
dragonfly_client_download_traffic{application="two",type="REMOTE_PEER"} 2e3
`,
			want: DownloadTraffic{RemotePeer: 3024},
		},
		{
			name:  "uninitialized metric family",
			input: "# TYPE unrelated counter\nunrelated 10\n",
		},
		{
			name:  "uninitialized remote peer series",
			input: counter + "dragonfly_client_download_traffic{type=\"BACK_TO_SOURCE\"} 7\n",
			want:  DownloadTraffic{BackToSource: 7},
		},
		{
			name:  "uninitialized back to source series",
			input: counter + "dragonfly_client_download_traffic{type=\"REMOTE_PEER\"} 7\n",
			want:  DownloadTraffic{RemotePeer: 7},
		},
		{
			name:  "malformed exposition",
			input: counter + "dragonfly_client_download_traffic{type=\"REMOTE_PEER\"} invalid\n", wantErr: true,
		},
		{
			name:  "nan",
			input: counter + "dragonfly_client_download_traffic{type=\"REMOTE_PEER\"} NaN\n", wantErr: true,
		},
		{
			name:  "infinite",
			input: counter + "dragonfly_client_download_traffic{type=\"REMOTE_PEER\"} +Inf\n", wantErr: true,
		},
		{
			name:  "negative",
			input: counter + "dragonfly_client_download_traffic{type=\"BACK_TO_SOURCE\"} -1\n", wantErr: true,
		},
		{
			name:  "fractional bytes",
			input: counter + "dragonfly_client_download_traffic{type=\"BACK_TO_SOURCE\"} 1.5\n", wantErr: true,
		},
		{
			name:  "overflowing sample",
			input: counter + "dragonfly_client_download_traffic{type=\"REMOTE_PEER\"} 18446744073709551616\n", wantErr: true,
		},
		{
			name: "overflowing aggregate",
			input: counter + `dragonfly_client_download_traffic{type="REMOTE_PEER",application="one"} 9223372036854775808
dragonfly_client_download_traffic{type="REMOTE_PEER",application="two"} 9223372036854775808
`,
			wantErr: true,
		},
		{
			name:  "missing type label",
			input: counter + "dragonfly_client_download_traffic 7\n", wantErr: true,
		},
		{
			name: "wrong metric type",
			input: "# TYPE dragonfly_client_download_traffic gauge\n" +
				"dragonfly_client_download_traffic{type=\"REMOTE_PEER\"} 7\n", wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDownloadTraffic(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseDownloadTraffic() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("parseDownloadTraffic() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDownloadTrafficDelta(t *testing.T) {
	baseline := DownloadTrafficSnapshot{
		PodUID: "pod-one", RestartCount: 1,
		Traffic: DownloadTraffic{RemotePeer: 10, BackToSource: 20},
	}
	tests := []struct {
		name    string
		before  DownloadTrafficSnapshot
		after   DownloadTrafficSnapshot
		want    DownloadTraffic
		wantErr bool
	}{
		{
			name: "traffic since baseline", before: baseline,
			after: DownloadTrafficSnapshot{
				PodUID: "pod-one", RestartCount: 1,
				Traffic: DownloadTraffic{RemotePeer: 33554442, BackToSource: 20},
			},
			want: DownloadTraffic{RemotePeer: 33554432},
		},
		{
			name: "unchanged counters", before: baseline, after: baseline,
		},
		{
			name: "successful source traffic is observable", before: baseline,
			after: DownloadTrafficSnapshot{
				PodUID: "pod-one", RestartCount: 1,
				Traffic: DownloadTraffic{RemotePeer: 10, BackToSource: 4116},
			},
			want: DownloadTraffic{BackToSource: 4096},
		},
		{
			name: "missing baseline UID", before: DownloadTrafficSnapshot{}, after: baseline, wantErr: true,
		},
		{
			name: "missing current UID", before: baseline, after: DownloadTrafficSnapshot{}, wantErr: true,
		},
		{
			name: "pod replaced despite increasing counters", before: baseline,
			after: DownloadTrafficSnapshot{
				PodUID: "pod-two", RestartCount: 1,
				Traffic: DownloadTraffic{RemotePeer: 100, BackToSource: 100},
			},
			wantErr: true,
		},
		{
			name: "container restarted despite increasing counters", before: baseline,
			after: DownloadTrafficSnapshot{
				PodUID: "pod-one", RestartCount: 2,
				Traffic: DownloadTraffic{RemotePeer: 100, BackToSource: 100},
			},
			wantErr: true,
		},
		{
			name: "restart count decreased", before: baseline,
			after: DownloadTrafficSnapshot{PodUID: "pod-one", Traffic: baseline.Traffic}, wantErr: true,
		},
		{
			name:   "negative restart count",
			before: DownloadTrafficSnapshot{PodUID: "pod-one", RestartCount: -1},
			after:  baseline, wantErr: true,
		},
		{
			name: "remote peer counter decreased", before: baseline,
			after: DownloadTrafficSnapshot{
				PodUID: "pod-one", RestartCount: 1,
				Traffic: DownloadTraffic{RemotePeer: 9, BackToSource: 20},
			},
			wantErr: true,
		},
		{
			name: "back to source counter decreased", before: baseline,
			after: DownloadTrafficSnapshot{
				PodUID: "pod-one", RestartCount: 1,
				Traffic: DownloadTraffic{RemotePeer: 10, BackToSource: 19},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadTrafficDelta(tt.before, tt.after)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DownloadTrafficDelta() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("DownloadTrafficDelta() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
