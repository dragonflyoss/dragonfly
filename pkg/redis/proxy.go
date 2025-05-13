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
	from string
	to   string
	done chan struct{}
}

// NewProxy creates a new proxy instance for redirecting traffic to redis.
func NewProxy(from string, to string) Proxy {
	return &proxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
	}
}

// Serve starts the proxy server and listens for incoming connections.
func (p *proxy) Serve() error {
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		select {
		case <-p.done:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				logger.Errorf("error accepting conn: %v", err)
			} else {
				go p.handleConn(conn)
			}
		}
	}
}

// Stop stops the proxy server and closes all connections.
func (p *proxy) Stop() {
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
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
	go p.copy(rConn, conn, wg)
	go p.copyReader(reader, rConn, wg)
	wg.Wait()
}

// copy copies data from one connection to another.
func (p *proxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			logger.Errorf("error copying from %s to %s: %v", from.RemoteAddr(), to.RemoteAddr(), err)
			p.Stop()
			return
		}
	}
}

// copyReader copies data from a reader to a connection.
func (p *proxy) copyReader(from io.Reader, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			logger.Errorf("error copying to %s: %v", to.RemoteAddr(), err)
			p.Stop()
			return
		}
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
