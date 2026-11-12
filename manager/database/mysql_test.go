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

package database

import (
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/manager/config"
)

func TestFormatMysqlDSN(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.MysqlConfig
		expect func(t *testing.T, dsn string, err error)
	}{
		{
			name: "plain connection without tls",
			cfg: &config.MysqlConfig{
				User:     "foo",
				Password: "bar",
				Host:     "localhost",
				Port:     3306,
				DBName:   "manager",
			},
			expect: func(t *testing.T, dsn string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				parsed, err := mysql.ParseDSN(dsn)
				assert.NoError(err)
				assert.Equal("foo", parsed.User)
				assert.Equal("bar", parsed.Passwd)
				assert.Equal("tcp", parsed.Net)
				assert.Equal("localhost:3306", parsed.Addr)
				assert.Equal("manager", parsed.DBName)
				assert.Empty(parsed.TLSConfig)
				assert.True(parsed.ParseTime)
				assert.True(parsed.InterpolateParams)
				assert.True(parsed.AllowNativePasswords)
				assert.Equal(time.Local, parsed.Loc)
				assert.Equal(defaultMysqlDialTimeout, parsed.Timeout)
				assert.Equal(defaultMysqlReadTimeout, parsed.ReadTimeout)
				assert.Equal(defaultMysqlWriteTimeout, parsed.WriteTimeout)
			},
		},
		{
			name: "ipv6 host is bracketed in the address",
			cfg: &config.MysqlConfig{
				User:     "foo",
				Password: "bar",
				Host:     "::1",
				Port:     3307,
				DBName:   "manager",
			},
			expect: func(t *testing.T, dsn string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				parsed, err := mysql.ParseDSN(dsn)
				assert.NoError(err)
				assert.Equal("[::1]:3307", parsed.Addr)
			},
		},
		{
			name: "tls config name is passed through when no custom tls is set",
			cfg: &config.MysqlConfig{
				User:      "foo",
				Password:  "bar",
				Host:      "localhost",
				Port:      3306,
				DBName:    "manager",
				TLSConfig: "preferred",
			},
			expect: func(t *testing.T, dsn string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				parsed, err := mysql.ParseDSN(dsn)
				assert.NoError(err)
				assert.Equal("preferred", parsed.TLSConfig)
			},
		},
		{
			name: "custom tls with insecure skip verify registers the custom config",
			cfg: &config.MysqlConfig{
				User:      "foo",
				Password:  "bar",
				Host:      "localhost",
				Port:      3306,
				DBName:    "manager",
				TLSConfig: "preferred",
				TLS:       &config.MysqlTLSClientConfig{InsecureSkipVerify: true},
			},
			expect: func(t *testing.T, dsn string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				parsed, err := mysql.ParseDSN(dsn)
				assert.NoError(err)
				assert.Equal("custom", parsed.TLSConfig)
				assert.NotNil(parsed.TLS)
				assert.True(parsed.TLS.InsecureSkipVerify)
			},
		},
		{
			name: "custom tls with unreadable ca cert fails",
			cfg: &config.MysqlConfig{
				User:     "foo",
				Password: "bar",
				Host:     "localhost",
				Port:     3306,
				DBName:   "manager",
				TLS:      &config.MysqlTLSClientConfig{CACert: "/nonexistent/ca.crt"},
			},
			expect: func(t *testing.T, dsn string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(dsn)
			},
		},
		{
			name: "custom tls with unreadable client key pair fails",
			cfg: &config.MysqlConfig{
				User:     "foo",
				Password: "bar",
				Host:     "localhost",
				Port:     3306,
				DBName:   "manager",
				TLS:      &config.MysqlTLSClientConfig{InsecureSkipVerify: true, Cert: "/nonexistent/client.crt", Key: "/nonexistent/client.key"},
			},
			expect: func(t *testing.T, dsn string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(dsn)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dsn, err := formatMysqlDSN(tc.cfg)
			tc.expect(t, dsn, err)
		})
	}
}
