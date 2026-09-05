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

package jwt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	// DefaultIssuer is the default issuer of inter-component gRPC JWTs.
	DefaultIssuer = "dragonfly-internal"

	// DefaultTokenTTL is the default lifetime of a generated JWT.
	DefaultTokenTTL = 10 * time.Minute

	// DefaultMaxTokenTTL is the maximum accepted lifetime of a JWT.
	DefaultMaxTokenTTL = 15 * time.Minute

	// DefaultClockSkew is the default clock skew allowed by the verifier.
	DefaultClockSkew = 30 * time.Second

	// DefaultRefreshBefore is the default duration before expiration at which a
	// cached JWT is refreshed.
	DefaultRefreshBefore = time.Minute

	// MinimumKeySize is the minimum decoded HMAC key size in bytes.
	MinimumKeySize = 32
)

// Mode controls client and server authentication behavior.
type Mode string

const (
	// ModeDisabled does not send or verify JWTs.
	ModeDisabled Mode = "disabled"

	// ModePermissive sends JWTs and accepts requests without JWTs during a
	// rolling upgrade. Invalid supplied JWTs are still rejected.
	ModePermissive Mode = "permissive"

	// ModeRequired sends JWTs and rejects requests without a valid JWT.
	ModeRequired Mode = "required"
)

// Config is the configuration for inter-component gRPC authentication.
type Config struct {
	// Mode controls client and server authentication behavior.
	Mode Mode `yaml:"mode" mapstructure:"mode"`

	// RequireTransportSecurity rejects sending JWTs over plaintext transports
	// when it is true.
	RequireTransportSecurity bool `yaml:"requireTransportSecurity" mapstructure:"requireTransportSecurity"`

	// JWT is the JWT configuration.
	JWT JWTConfig `yaml:"jwt" mapstructure:"jwt"`
}

// JWTConfig is the configuration for signing and verifying JWTs.
type JWTConfig struct {
	// Issuer is the exact issuer required by the verifier.
	Issuer string `yaml:"issuer" mapstructure:"issuer"`

	// TokenTTL is the lifetime of a generated JWT.
	TokenTTL time.Duration `yaml:"tokenTTL" mapstructure:"tokenTTL"`

	// MaxTokenTTL is the maximum accepted lifetime of a JWT.
	MaxTokenTTL time.Duration `yaml:"maxTokenTTL" mapstructure:"maxTokenTTL"`

	// ClockSkew is the allowed clock difference between components.
	ClockSkew time.Duration `yaml:"clockSkew" mapstructure:"clockSkew"`

	// RefreshBefore controls when a cached JWT is refreshed.
	RefreshBefore time.Duration `yaml:"refreshBefore" mapstructure:"refreshBefore"`

	// ActiveKeyID identifies the key used to sign new JWTs.
	ActiveKeyID string `yaml:"activeKeyID" mapstructure:"activeKeyID"`

	// Keys contains all keys trusted by the verifier.
	Keys []KeyConfig `yaml:"keys" mapstructure:"keys"`
}

// KeyConfig identifies a shared HMAC key file.
type KeyConfig struct {
	// ID is serialized as the JWT kid header.
	ID string `yaml:"id" mapstructure:"id"`

	// SecretFile contains a Base64-encoded HMAC key.
	SecretFile string `yaml:"secretFile" mapstructure:"secretFile"`
}

// DefaultConfig returns the default disabled gRPC authentication configuration.
func DefaultConfig() Config {
	return Config{
		Mode: ModeDisabled,
		JWT: JWTConfig{
			Issuer:        DefaultIssuer,
			TokenTTL:      DefaultTokenTTL,
			MaxTokenTTL:   DefaultMaxTokenTTL,
			ClockSkew:     DefaultClockSkew,
			RefreshBefore: DefaultRefreshBefore,
		},
	}
}

// EffectiveMode treats an omitted mode as disabled for backward compatibility.
func (c Config) EffectiveMode() Mode {
	if c.Mode == "" {
		return ModeDisabled
	}

	return c.Mode
}

// Enabled returns whether gRPC authentication is enabled.
func (c Config) Enabled() bool {
	return c.EffectiveMode() != ModeDisabled
}

// ValidateConfig validates configuration and referenced key files.
func ValidateConfig(config Config) error {
	_, err := validateAndLoadKeys(config)
	return err
}

func validateAndLoadKeys(config Config) (map[string][]byte, error) {
	switch config.EffectiveMode() {
	case ModeDisabled:
		return nil, nil
	case ModePermissive, ModeRequired:
	default:
		return nil, fmt.Errorf("grpc auth has unsupported mode %q", config.Mode)
	}

	if config.JWT.Issuer == "" {
		return nil, errors.New("grpc auth jwt requires parameter issuer")
	}

	if config.JWT.TokenTTL < time.Second {
		return nil, errors.New("grpc auth jwt tokenTTL must be at least one second")
	}

	if config.JWT.MaxTokenTTL < time.Second {
		return nil, errors.New("grpc auth jwt maxTokenTTL must be at least one second")
	}

	if config.JWT.TokenTTL%time.Second != 0 ||
		config.JWT.MaxTokenTTL%time.Second != 0 ||
		config.JWT.ClockSkew%time.Second != 0 ||
		config.JWT.RefreshBefore%time.Second != 0 {
		return nil, errors.New("grpc auth jwt durations must use whole seconds")
	}

	if config.JWT.TokenTTL > config.JWT.MaxTokenTTL {
		return nil, errors.New("grpc auth jwt tokenTTL must not exceed maxTokenTTL")
	}

	if config.JWT.ClockSkew < 0 {
		return nil, errors.New("grpc auth jwt clockSkew must not be negative")
	}

	if config.JWT.RefreshBefore <= 0 || config.JWT.RefreshBefore >= config.JWT.TokenTTL {
		return nil, errors.New("grpc auth jwt refreshBefore must be positive and less than tokenTTL")
	}

	if config.JWT.ActiveKeyID == "" {
		return nil, errors.New("grpc auth jwt requires parameter activeKeyID")
	}

	if len(config.JWT.Keys) == 0 {
		return nil, errors.New("grpc auth jwt requires at least one key")
	}

	keys := make(map[string][]byte, len(config.JWT.Keys))
	for _, keyConfig := range config.JWT.Keys {
		if keyConfig.ID == "" {
			return nil, errors.New("grpc auth jwt key requires parameter id")
		}

		if _, ok := keys[keyConfig.ID]; ok {
			return nil, fmt.Errorf("grpc auth jwt key id %q is duplicated", keyConfig.ID)
		}

		if keyConfig.SecretFile == "" {
			return nil, fmt.Errorf("grpc auth jwt key %q requires parameter secretFile", keyConfig.ID)
		}

		secret, err := loadKey(keyConfig.SecretFile)
		if err != nil {
			return nil, fmt.Errorf("grpc auth jwt key %q: %w", keyConfig.ID, err)
		}

		keys[keyConfig.ID] = secret
	}

	if _, ok := keys[config.JWT.ActiveKeyID]; !ok {
		return nil, fmt.Errorf("grpc auth jwt active key id %q is not trusted", config.JWT.ActiveKeyID)
	}

	return keys, nil
}

func loadKey(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat secret file: %w", err)
	}

	if !info.Mode().IsRegular() {
		return nil, errors.New("secret file is not a regular file")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read secret file: %w", err)
	}

	secret, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(content)))
	if err != nil {
		return nil, errors.New("secret file is not valid standard Base64")
	}

	if len(secret) < MinimumKeySize {
		return nil, fmt.Errorf("decoded secret must contain at least %d bytes", MinimumKeySize)
	}

	return secret, nil
}
