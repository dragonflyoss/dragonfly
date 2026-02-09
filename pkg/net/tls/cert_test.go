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
	"reflect"
	"testing"
	"time"
)

func TestPEMToCertPool(t *testing.T) {
	// Generate a valid certificate for testing.
	_, pemCert, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("failed to generate test certificate: %v", err)
	}

	// Create an expected CertPool with the test cert.
	expectedCertPool := x509.NewCertPool()
	if !expectedCertPool.AppendCertsFromPEM(pemCert) {
		t.Fatal("failed to create expected cert pool")
	}

	tests := []struct {
		name        string
		pemCerts    []byte
		expected    *x509.CertPool
		expectedErr bool
		checkPool   bool // Flag to check if the pool content should be compared
	}{
		{
			name:        "Empty PEM certs",
			pemCerts:    []byte{},
			expected:    nil,
			expectedErr: false,
		},
		{
			name:        "Valid PEM cert",
			pemCerts:    pemCert,
			expected:    expectedCertPool,
			expectedErr: false,
			checkPool:   true,
		},
		{
			name:        "Invalid PEM cert",
			pemCerts:    []byte("this is not a valid pem"),
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PEMToCertPool(tt.pemCerts)
			if (err != nil) != tt.expectedErr {
				t.Errorf("PEMToCertPool() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if tt.checkPool {
				if got == nil || len(got.Subjects()) != len(tt.expected.Subjects()) {
					t.Errorf("PEMToCertPool() = %v, expected %v", got, tt.expected)
				}
			} else {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("PEMToCertPool() = %v, expected %v", got, tt.expected)
				}
			}
		})
	}
}

func TestDERToCertPool(t *testing.T) {
	// Generate a valid certificate for testing.
	derCert, _, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("failed to generate test certificate: %v", err)
	}

	// Create an expected CertPool with the test cert.
	cert, err := x509.ParseCertificate(derCert)
	if err != nil {
		t.Fatalf("failed to parse generated der cert: %v", err)
	}
	expectedCertPool := x509.NewCertPool()
	expectedCertPool.AddCert(cert)

	tests := []struct {
		name        string
		derCerts    [][]byte
		expected    *x509.CertPool
		expectedErr bool
		checkPool   bool // Flag to check if the pool content should be compared
	}{
		{
			name:        "Empty DER certs",
			derCerts:    [][]byte{},
			expected:    nil,
			expectedErr: false,
		},
		{
			name:        "Valid DER cert",
			derCerts:    [][]byte{derCert},
			expected:    expectedCertPool,
			expectedErr: false,
			checkPool:   true,
		},
		{
			name:        "Invalid DER cert",
			derCerts:    [][]byte{[]byte("this is not a valid der")},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DERToCertPool(tt.derCerts)
			if (err != nil) != tt.expectedErr {
				t.Errorf("DERToCertPool() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if tt.checkPool {
				if got == nil || len(got.Subjects()) != len(tt.expected.Subjects()) {
					t.Errorf("DERToCertPool() = %v, expected %v", got, tt.expected)
				}
			} else {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("DERToCertPool() = %v, expected %v", got, tt.expected)
				}
			}
		})
	}
}

// generateTestCertificate generates a self-signed certificate for testing purposes
// and returns the certificate in both DER and PEM formats.
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
