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

package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"d7y.io/dragonfly/v2/cmd/dependency/base"
	"d7y.io/dragonfly/v2/pkg/types"
)

const (
	// DefaultLeaderElectionID is the name of the lease used for leader election.
	DefaultLeaderElectionID = "datacontroller.data.d7y.io"

	// DefaultMetricsAddr is the address the metrics endpoint listens on.
	DefaultMetricsAddr = ":8000"

	// DefaultHealthProbeAddr is the address the health probe endpoint listens on.
	DefaultHealthProbeAddr = ":8001"

	// DefaultPollInterval is how often in-flight manager jobs are polled.
	DefaultPollInterval = 15 * time.Second

	// DefaultResyncInterval is how often every object is reconciled without an event.
	DefaultResyncInterval = 10 * time.Minute

	// DefaultMaxConcurrentReconciles is the number of datasets reconciled in parallel.
	DefaultMaxConcurrentReconciles = 4

	// DefaultManagerRequestTimeout bounds every request to the manager.
	DefaultManagerRequestTimeout = 30 * time.Second

	// DefaultRetryBackoff is the initial delay before a failed operation is retried.
	DefaultRetryBackoff = 30 * time.Second

	// DefaultMaxRetryBackoff caps the delay before a failed operation is retried.
	DefaultMaxRetryBackoff = 10 * time.Minute
)

// Config is the configuration of the data controller.
type Config struct {
	// Base options.
	base.Options `yaml:",inline" mapstructure:",squash"`

	// Server configuration.
	Server ServerConfig `yaml:"server" mapstructure:"server"`

	// Manager configuration.
	Manager ManagerConfig `yaml:"manager" mapstructure:"manager"`

	// Controller configuration.
	Controller ControllerConfig `yaml:"controller" mapstructure:"controller"`

	// Metrics configuration.
	Metrics MetricsConfig `yaml:"metrics" mapstructure:"metrics"`
}

// ServerConfig configures logging.
type ServerConfig struct {
	// LogDir is the log directory.
	LogDir string `yaml:"logDir" mapstructure:"logDir"`

	// LogLevel is the log level, one of "debug", "info", "warn", "error", "panic" and "fatal".
	LogLevel string `yaml:"logLevel" mapstructure:"logLevel"`

	// LogMaxSize is the maximum size in megabytes of a log file before rotation.
	LogMaxSize int `yaml:"logMaxSize" mapstructure:"logMaxSize"`

	// LogMaxAge is the maximum number of days to retain old log files.
	LogMaxAge int `yaml:"logMaxAge" mapstructure:"logMaxAge"`

	// LogMaxBackups is the maximum number of old log files to keep.
	LogMaxBackups int `yaml:"logMaxBackups" mapstructure:"logMaxBackups"`
}

// ManagerConfig configures the default Dragonfly manager. A Dataset can override
// it with spec.managerRef.
type ManagerConfig struct {
	// Endpoint is the base URL of the manager REST API, for example http://dragonfly-manager:8080.
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`

	// Token is a personal access token with the job scope. Prefer TokenFile or the
	// DATACONTROLLER_MANAGER_TOKEN environment variable over writing it to the config file.
	Token string `yaml:"token" mapstructure:"token"`

	// TokenFile is a file holding the personal access token, for example a mounted Secret.
	TokenFile string `yaml:"tokenFile" mapstructure:"tokenFile"`

	// InsecureSkipTLSVerify disables TLS verification against the manager.
	InsecureSkipTLSVerify bool `yaml:"insecureSkipTLSVerify" mapstructure:"insecureSkipTLSVerify"`

	// RequestTimeout bounds every request to the manager.
	RequestTimeout time.Duration `yaml:"requestTimeout" mapstructure:"requestTimeout"`
}

// ControllerConfig configures the reconcilers.
type ControllerConfig struct {
	// KubeConfig is the path of a kubeconfig file. Empty uses the in-cluster configuration.
	KubeConfig string `yaml:"kubeConfig" mapstructure:"kubeConfig"`

	// Namespace restricts the controller to a single namespace. Empty watches every namespace.
	Namespace string `yaml:"namespace" mapstructure:"namespace"`

	// LeaderElection enables leader election so several replicas can run.
	LeaderElection bool `yaml:"leaderElection" mapstructure:"leaderElection"`

	// LeaderElectionNamespace is the namespace of the leader election lease. Empty uses
	// the namespace the controller runs in.
	LeaderElectionNamespace string `yaml:"leaderElectionNamespace" mapstructure:"leaderElectionNamespace"`

	// LeaderElectionID is the name of the leader election lease.
	LeaderElectionID string `yaml:"leaderElectionID" mapstructure:"leaderElectionID"`

	// HealthProbeAddr is the address the liveness and readiness endpoints listen on.
	HealthProbeAddr string `yaml:"healthProbeAddr" mapstructure:"healthProbeAddr"`

	// MaxConcurrentReconciles is the number of datasets reconciled in parallel.
	MaxConcurrentReconciles int `yaml:"maxConcurrentReconciles" mapstructure:"maxConcurrentReconciles"`

	// PollInterval is how often in-flight manager jobs are polled.
	PollInterval time.Duration `yaml:"pollInterval" mapstructure:"pollInterval"`

	// ResyncInterval is how often every object is reconciled without an event.
	ResyncInterval time.Duration `yaml:"resyncInterval" mapstructure:"resyncInterval"`

	// RetryBackoff is the initial delay before a failed operation is retried; it doubles on
	// every consecutive failure.
	RetryBackoff time.Duration `yaml:"retryBackoff" mapstructure:"retryBackoff"`

	// MaxRetryBackoff caps the delay before a failed operation is retried.
	MaxRetryBackoff time.Duration `yaml:"maxRetryBackoff" mapstructure:"maxRetryBackoff"`
}

// MetricsConfig configures the metrics endpoint.
type MetricsConfig struct {
	// Enable the metrics endpoint.
	Enable bool `yaml:"enable" mapstructure:"enable"`

	// Addr is the address the metrics endpoint listens on.
	Addr string `yaml:"addr" mapstructure:"addr"`
}

// New returns a Config with defaults.
func New() *Config {
	return &Config{
		Options: base.Options{
			Tracing: base.TracingConfig{
				ServiceName: types.DataControllerName,
			},
		},
		Server: ServerConfig{
			LogLevel:      "info",
			LogMaxSize:    1024,
			LogMaxAge:     7,
			LogMaxBackups: 20,
		},
		Manager: ManagerConfig{
			RequestTimeout: DefaultManagerRequestTimeout,
		},
		Controller: ControllerConfig{
			LeaderElectionID:        DefaultLeaderElectionID,
			HealthProbeAddr:         DefaultHealthProbeAddr,
			MaxConcurrentReconciles: DefaultMaxConcurrentReconciles,
			PollInterval:            DefaultPollInterval,
			ResyncInterval:          DefaultResyncInterval,
			RetryBackoff:            DefaultRetryBackoff,
			MaxRetryBackoff:         DefaultMaxRetryBackoff,
		},
		Metrics: MetricsConfig{
			Enable: true,
			Addr:   DefaultMetricsAddr,
		},
	}
}

// Validate checks the configuration.
func (cfg *Config) Validate() error {
	if cfg.Server.LogLevel != "" {
		switch strings.ToLower(cfg.Server.LogLevel) {
		case "debug", "info", "warn", "error", "panic", "fatal":
		default:
			return fmt.Errorf("server requires parameter logLevel to be one of debug, info, warn, error, panic, fatal, got %q", cfg.Server.LogLevel)
		}
	}

	if cfg.Manager.Endpoint != "" {
		u, err := url.Parse(cfg.Manager.Endpoint)
		if err != nil {
			return fmt.Errorf("manager requires parameter endpoint to be a URL: %w", err)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New("manager requires parameter endpoint to use http or https")
		}

		if cfg.Manager.Token == "" && cfg.Manager.TokenFile == "" {
			return errors.New("manager requires parameter token or tokenFile")
		}
	}

	if cfg.Manager.RequestTimeout <= 0 {
		return errors.New("manager requires parameter requestTimeout")
	}

	if cfg.Controller.LeaderElection && cfg.Controller.LeaderElectionID == "" {
		return errors.New("controller requires parameter leaderElectionID")
	}

	if cfg.Controller.MaxConcurrentReconciles <= 0 {
		return errors.New("controller requires parameter maxConcurrentReconciles")
	}

	if cfg.Controller.PollInterval <= 0 {
		return errors.New("controller requires parameter pollInterval")
	}

	if cfg.Controller.ResyncInterval <= 0 {
		return errors.New("controller requires parameter resyncInterval")
	}

	if cfg.Controller.RetryBackoff <= 0 {
		return errors.New("controller requires parameter retryBackoff")
	}

	if cfg.Controller.MaxRetryBackoff < cfg.Controller.RetryBackoff {
		return errors.New("controller requires parameter maxRetryBackoff to be greater than or equal to retryBackoff")
	}

	if cfg.Metrics.Enable && cfg.Metrics.Addr == "" {
		return errors.New("metrics requires parameter addr")
	}

	return nil
}

// Convert resolves values that are derived from other values, such as the token file.
func (cfg *Config) Convert() error {
	if cfg.Manager.Token == "" && cfg.Manager.TokenFile != "" {
		data, err := os.ReadFile(cfg.Manager.TokenFile)
		if err != nil {
			return fmt.Errorf("read manager token file %q: %w", cfg.Manager.TokenFile, err)
		}

		cfg.Manager.Token = strings.TrimSpace(string(data))
		if cfg.Manager.Token == "" {
			return fmt.Errorf("manager token file %q is empty", cfg.Manager.TokenFile)
		}
	}

	return nil
}
