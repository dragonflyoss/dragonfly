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
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials"
)

type mockCertFiles struct {
	caCertFile string
	certFile   string
	keyFile    string
}

func mockCertificates(t *testing.T, dir string) mockCertFiles {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Test Co"}},
		DNSNames:              []string{"localhost"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	files := mockCertFiles{
		caCertFile: filepath.Join(dir, "ca.pem"),
		certFile:   filepath.Join(dir, "cert.pem"),
		keyFile:    filepath.Join(dir, "key.pem"),
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err := os.WriteFile(files.caCertFile, certPEM, 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(files.certFile, certPEM, 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(files.keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}), 0600); err != nil {
		t.Fatal(err)
	}

	return files
}

func TestNewServerCredentials(t *testing.T) {
	dir := t.TempDir()
	files := mockCertificates(t, dir)
	invalidCAFile := filepath.Join(dir, "invalid-ca.pem")
	if err := os.WriteFile(invalidCAFile, []byte("this is not a valid pem"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		caCertFile string
		certFile   string
		keyFile    string
		expect     func(t *testing.T, creds credentials.TransportCredentials, err error)
	}{
		{
			name:       "valid ca certificate and key pair",
			caCertFile: files.caCertFile,
			certFile:   files.certFile,
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("tls", creds.Info().SecurityProtocol)
			},
		},
		{
			name:       "missing certificate file",
			caCertFile: files.caCertFile,
			certFile:   filepath.Join(dir, "missing-cert.pem"),
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
		{
			name:       "missing key file",
			caCertFile: files.caCertFile,
			certFile:   files.certFile,
			keyFile:    filepath.Join(dir, "missing-key.pem"),
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
		{
			name:       "missing ca certificate file",
			caCertFile: filepath.Join(dir, "missing-ca.pem"),
			certFile:   files.certFile,
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
		{
			name:       "ca certificate file without valid pem",
			caCertFile: invalidCAFile,
			certFile:   files.certFile,
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			creds, err := NewServerCredentials(tc.caCertFile, tc.certFile, tc.keyFile)
			tc.expect(t, creds, err)
		})
	}
}

func TestNewClientCredentials(t *testing.T) {
	dir := t.TempDir()
	files := mockCertificates(t, dir)
	invalidCAFile := filepath.Join(dir, "invalid-ca.pem")
	if err := os.WriteFile(invalidCAFile, []byte("this is not a valid pem"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		caCertFile string
		certFile   string
		keyFile    string
		expect     func(t *testing.T, creds credentials.TransportCredentials, err error)
	}{
		{
			name:       "valid ca certificate and key pair",
			caCertFile: files.caCertFile,
			certFile:   files.certFile,
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("tls", creds.Info().SecurityProtocol)
			},
		},
		{
			name:       "missing certificate file",
			caCertFile: files.caCertFile,
			certFile:   filepath.Join(dir, "missing-cert.pem"),
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
		{
			name:       "missing key file",
			caCertFile: files.caCertFile,
			certFile:   files.certFile,
			keyFile:    filepath.Join(dir, "missing-key.pem"),
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
		{
			name:       "missing ca certificate file",
			caCertFile: filepath.Join(dir, "missing-ca.pem"),
			certFile:   files.certFile,
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
		{
			name:       "ca certificate file without valid pem",
			caCertFile: invalidCAFile,
			certFile:   files.certFile,
			keyFile:    files.keyFile,
			expect: func(t *testing.T, creds credentials.TransportCredentials, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(creds)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			creds, err := NewClientCredentials(tc.caCertFile, tc.certFile, tc.keyFile)
			tc.expect(t, creds, err)
		})
	}
}

func TestCredentials_Handshake(t *testing.T) {
	server := mockCertificates(t, t.TempDir())
	other := mockCertificates(t, t.TempDir())
	tests := []struct {
		name         string
		clientCAFile string
		expect       func(t *testing.T, clientErr, serverErr error)
	}{
		{
			name:         "mutual tls with shared ca",
			clientCAFile: server.caCertFile,
			expect: func(t *testing.T, clientErr, serverErr error) {
				assert := assert.New(t)
				assert.NoError(clientErr)
				assert.NoError(serverErr)
			},
		},
		{
			name:         "client rejects server certificate signed by untrusted ca",
			clientCAFile: other.caCertFile,
			expect: func(t *testing.T, clientErr, serverErr error) {
				assert := assert.New(t)
				assert.Error(clientErr)
				assert.Error(serverErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			serverCreds, err := NewServerCredentials(server.caCertFile, server.certFile, server.keyFile)
			assert.NoError(err)
			clientCreds, err := NewClientCredentials(tc.clientCAFile, server.certFile, server.keyFile)
			assert.NoError(err)

			clientConn, serverConn := net.Pipe()
			defer clientConn.Close()
			defer serverConn.Close()
			serverErrCh := make(chan error, 1)
			go func() {
				_, _, err := serverCreds.ServerHandshake(serverConn)
				serverErrCh <- err
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, _, clientErr := clientCreds.ClientHandshake(ctx, "localhost", clientConn)
			tc.expect(t, clientErr, <-serverErrCh)
		})
	}
}
