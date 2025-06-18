/*
 *     Copyright 2020 The Dragonfly Authors
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

package base

type Options struct {
	Console   bool          `yaml:"console" mapstructure:"console"`
	Verbose   bool          `yaml:"verbose" mapstructure:"verbose"`
	PProfPort int           `yaml:"pprof-port" mapstructure:"pprof-port"`
	Tracing   TracingConfig `yaml:"tracing" mapstructure:"tracing"`
}

// TracingConfig defines the configuration for OpenTelemetry tracing.
type TracingConfig struct {
	// Addr is the address of the tracing collector.
	Addr string `yaml:"addr" mapstructure:"addr"`
	// ServiceName is the name of the service for tracing.
	ServiceName string `yaml:"service-name" mapstructure:"service-name"`
	// Headers are additional headers to be sent with tracing requests.
	Headers map[string]string `yaml:"headers" mapstructure:"headers"`
}
