/*
 *     Copyright 2025 The Dragonfly Authors
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

package redis

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const mockRedisRequest = "*1\r\n$4\r\nPING\r\n"

func mockFreeAddr(t *testing.T) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	return listener.Addr().String()
}

func mockDial(t *testing.T, addr string) net.Conn {
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			return conn
		}

		if time.Now().After(deadline) {
			t.Fatalf("dial %s: %v", addr, err)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func mockServeEchoBackend(listener net.Listener, accepted *atomic.Int32) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		accepted.Add(1)
		go func() {
			defer conn.Close()
			buf := make([]byte, len(mockRedisRequest))
			n, _ := io.ReadFull(conn, buf)
			_, _ = conn.Write(buf[:n])
		}()
	}
}

func TestProxy_Serve(t *testing.T) {
	tests := []struct {
		name    string
		request string
		expect  func(t *testing.T, response []byte, err error, backendConns int32)
	}{
		{
			name:    "redis protocol request is forwarded to the backend and its response returned",
			request: mockRedisRequest,
			expect: func(t *testing.T, response []byte, err error, backendConns int32) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(mockRedisRequest, string(response))
				assert.Equal(int32(1), backendConns)
			},
		},
		{
			name:    "non redis protocol request is closed without dialing the backend",
			request: "GET / HTTP/1.1\r\n",
			expect: func(t *testing.T, response []byte, err error, backendConns int32) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(response)
				assert.Equal(int32(0), backendConns)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			backend, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}

			defer backend.Close()

			var backendConns atomic.Int32
			go mockServeEchoBackend(backend, &backendConns)

			from := mockFreeAddr(t)
			p := NewProxy(from, backend.Addr().String())
			served := make(chan error, 1)
			go func() { served <- p.Serve() }()

			conn := mockDial(t, from)
			defer conn.Close()
			if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}

			if _, err := conn.Write([]byte(tc.request)); err != nil {
				t.Fatal(err)
			}

			response := make([]byte, len(tc.request))
			n, err := io.ReadFull(conn, response)
			tc.expect(t, response[:n], err, backendConns.Load())

			p.Stop()
			mockDial(t, from).Close()
			select {
			case err := <-served:
				assert.NoError(err)
			case <-time.After(5 * time.Second):
				t.Fatal("proxy did not stop serving after Stop and a wake-up connection")
			}
		})
	}
}

func TestProxy_isRedisProtocol(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect func(t *testing.T, ok bool, payload string)
	}{
		{
			name:  "array keeps the peeked bytes readable",
			input: mockRedisRequest,
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(mockRedisRequest, payload)
			},
		},
		{
			name:  "simple string",
			input: "+OK\r\n",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal("+OK\r\n", payload)
			},
		},
		{
			name:  "error",
			input: "-ERR unknown command\r\n",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal("-ERR unknown command\r\n", payload)
			},
		},
		{
			name:  "integer",
			input: ":1000\r\n",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal(":1000\r\n", payload)
			},
		},
		{
			name:  "bulk string",
			input: "$4\r\nPING\r\n",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.True(ok)
				assert.Equal("$4\r\nPING\r\n", payload)
			},
		},
		{
			name:  "inline command",
			input: "PING\r\n",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.False(ok)
				assert.Equal("PING\r\n", payload)
			},
		},
		{
			name:  "http request",
			input: "GET / HTTP/1.1\r\n",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.False(ok)
				assert.Equal("GET / HTTP/1.1\r\n", payload)
			},
		},
		{
			name:  "connection closed before any byte",
			input: "",
			expect: func(t *testing.T, ok bool, payload string) {
				assert := assert.New(t)
				assert.False(ok)
				assert.Empty(payload)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, server := net.Pipe()
			go func() {
				defer client.Close()
				if tc.input != "" {
					_, _ = client.Write([]byte(tc.input))
				}
			}()

			p := &proxy{}
			reader, ok := p.isRedisProtocol(server)
			payload, _ := io.ReadAll(reader)
			tc.expect(t, ok, string(payload))
		})
	}
}

func TestProxy_StopConcurrent(t *testing.T) {
	tests := []struct {
		name    string
		callers int
		expect  func(t *testing.T, stopAll func())
	}{
		{
			name:    "stop is idempotent when called once",
			callers: 1,
			expect: func(t *testing.T, stopAll func()) {
				assert := assert.New(t)
				assert.NotPanics(stopAll)
			},
		},
		{
			name:    "stop is idempotent under concurrent callers",
			callers: 16,
			expect: func(t *testing.T, stopAll func()) {
				assert := assert.New(t)
				assert.NotPanics(stopAll)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewProxy("127.0.0.1:0", "127.0.0.1:0")

			var wg sync.WaitGroup
			wg.Add(tc.callers)
			for i := 0; i < tc.callers; i++ {
				go func() {
					defer wg.Done()
					p.Stop()
				}()
			}

			tc.expect(t, wg.Wait)
		})
	}
}
