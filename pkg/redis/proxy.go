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
	"bufio"
	"errors"
	"io"
	"net"
	"sync"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

type Proxy interface {
	Serve() error
	Stop()
}

type proxy struct {
	from      string
	to        string
	done      chan struct{}
	closeOnce sync.Once
	mu        sync.Mutex
	listener  net.Listener
}

// NewProxy creates a new proxy instance for redirecting traffic to redis.
func NewProxy(from string, to string) Proxy {
	return &proxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
	}
}

// Serve starts the proxy server and listens for incoming connections until Stop is called.
func (p *proxy) Serve() error {
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	defer listener.Close()

	p.mu.Lock()
	select {
	case <-p.done:
		p.mu.Unlock()
		return nil
	default:
	}
	p.listener = listener
	p.mu.Unlock()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-p.done:
				return nil
			default:
			}

			if errors.Is(err, net.ErrClosed) {
				return err
			}

			logger.Errorf("error accepting conn: %v", err)
			continue
		}

		go p.handleConn(conn)
	}
}

// Stop stops the proxy server and unblocks Serve by closing the listener.
func (p *proxy) Stop() {
	p.closeOnce.Do(func() {
		close(p.done)

		p.mu.Lock()
		defer p.mu.Unlock()
		if p.listener != nil {
			p.listener.Close()
		}
	})
}

// handleConn handles the incoming connection and establishes a connection to the remote host.
func (p *proxy) handleConn(conn net.Conn) {
	defer conn.Close()

	reader, isRedisProtocol := p.isRedisProtocol(conn)
	if !isRedisProtocol {
		logger.Errorf("not a redis protocol: %s", conn.RemoteAddr())
		return
	}

	rConn, err := net.Dial("tcp", p.to)
	if err != nil {
		logger.Errorf("error dialing remote host: %v", err)
		return
	}
	defer rConn.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(conn, rConn, wg)
	go p.copy(rConn, reader, wg)
	wg.Wait()
}

// copy copies data from src to dst and closes dst when src is drained or fails,
// so the opposite direction of the same connection unblocks. A failure only
// affects this connection; the proxy keeps serving others.
func (p *proxy) copy(dst net.Conn, src io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil && !errors.Is(err, net.ErrClosed) {
		logger.Errorf("error copying to %s: %v", dst.RemoteAddr(), err)
	}
}

// isRedisProtocol checks if the connection uses the Redis protocol.
func (p *proxy) isRedisProtocol(conn net.Conn) (io.Reader, bool) {
	reader := bufio.NewReader(conn)
	firstByte, err := reader.Peek(1)
	if err != nil {
		if err != io.EOF {
			logger.Errorf("reading first byte from client failed: %s: %v", conn.RemoteAddr(), err)
		}

		return reader, false
	}

	switch firstByte[0] {
	case '*', '+', '-', ':', '$':
		return reader, true
	}

	return reader, false
}
