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

package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPEMToCertPool(t *testing.T) {
	_, pemCert, err := generateTestCertificate()
	if err != nil {
		t.Fatal(err)
	}

	expectedCertPool := x509.NewCertPool()
	if !expectedCertPool.AppendCertsFromPEM(pemCert) {
		t.FailNow()
	}

	tests := []struct {
		name     string
		pemCerts []byte
		expect   func(t *testing.T, pool *x509.CertPool, err error)
	}{
		{
			name:     "empty pem certs",
			pemCerts: []byte{},
			expect: func(t *testing.T, pool *x509.CertPool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Nil(pool)
			},
		},
		{
			name:     "valid pem cert",
			pemCerts: pemCert,
			expect: func(t *testing.T, pool *x509.CertPool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(expectedCertPool.Equal(pool))
			},
		},
		{
			name:     "invalid pem cert",
			pemCerts: []byte("this is not a valid pem"),
			expect: func(t *testing.T, pool *x509.CertPool, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(pool)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool, err := PEMToCertPool(tc.pemCerts)
			tc.expect(t, pool, err)
		})
	}
}

func TestDERToCertPool(t *testing.T) {
	derCert, _, err := generateTestCertificate()
	if err != nil {
		t.Fatal(err)
	}

	cert, err := x509.ParseCertificate(derCert)
	if err != nil {
		t.Fatal(err)
	}

	expectedCertPool := x509.NewCertPool()
	expectedCertPool.AddCert(cert)

	tests := []struct {
		name     string
		derCerts [][]byte
		expect   func(t *testing.T, pool *x509.CertPool, err error)
	}{
		{
			name:     "empty der certs",
			derCerts: [][]byte{},
			expect: func(t *testing.T, pool *x509.CertPool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Nil(pool)
			},
		},
		{
			name:     "valid der cert",
			derCerts: [][]byte{derCert},
			expect: func(t *testing.T, pool *x509.CertPool, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(expectedCertPool.Equal(pool))
			},
		},
		{
			name:     "invalid der cert",
			derCerts: [][]byte{[]byte("this is not a valid der")},
			expect: func(t *testing.T, pool *x509.CertPool, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(pool)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool, err := DERToCertPool(tc.derCerts)
			tc.expect(t, pool, err)
		})
	}
}

func generateTestCertificate() (derBytes []byte, pemBytes []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	pemBytes = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	return derBytes, pemBytes, nil
}
