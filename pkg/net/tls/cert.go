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
	"crypto/x509"
	"errors"
)

// PEMToCertPool concerts PEM certificates to a CertPool.
func PEMToCertPool(pemCerts []byte) (*x509.CertPool, error) {
	if len(pemCerts) == 0 {
		return nil, nil
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCerts) {
		return nil, errors.New("failed to load cert pool")
	}

	return certPool, nil
}

// DERToCertPool converts DER encoded certificates to a CertPool.
func DERToCertPool(derCerts [][]byte) (*x509.CertPool, error) {
	if len(derCerts) == 0 {
		return nil, nil
	}

	certPool := x509.NewCertPool()
	for _, derCert := range derCerts {
		cert, err := x509.ParseCertificate(derCert)
		if err != nil {
			return nil, err
		}

		certPool.AddCert(cert)
	}

	return certPool, nil
}
