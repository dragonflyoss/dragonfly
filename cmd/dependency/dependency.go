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

package dependency

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-echarts/statsview"
	"github.com/go-echarts/statsview/viewer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	"d7y.io/dragonfly/v2/cmd/dependency/base"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/dfnet"
	"d7y.io/dragonfly/v2/pkg/dfpath"
	"d7y.io/dragonfly/v2/pkg/net/fqdn"
	"d7y.io/dragonfly/v2/pkg/net/ip"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/version"
)

// InitCommandAndConfig initializes flags binding and common sub cmds.
// config is a pointer to configuration struct.
func InitCommandAndConfig(cmd *cobra.Command, useConfigFile bool, config any) {
	rootName := cmd.Root().Name()
	cobra.OnInitialize(func() {
		initConfig(useConfigFile, rootName, config)
	})

	if !cmd.HasParent() {
		// Add common flags
		flags := cmd.PersistentFlags()
		flags.Bool("console", false, "whether logger output records to the stdout")
		flags.String("config", "", fmt.Sprintf("the path of configuration file with yaml extension name, default is %s, it can also be set by env var: %s", filepath.Join(dfpath.DefaultConfigDir, rootName+".yaml"), strings.ToUpper(rootName+"_config")))

		// Bind common flags
		if err := viper.BindPFlags(flags); err != nil {
			panic(fmt.Errorf("bind common flags to viper: %w", err))
		}

		// Config for binding env
		viper.SetEnvPrefix(rootName)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		_ = viper.BindEnv("config")

		// Add common cmds only on root cmd
		cmd.AddCommand(VersionCmd)
		cmd.AddCommand(newDocCommand(cmd.Name()))
		cmd.AddCommand(PluginCmd)
	}
}

// InitMonitor initialize monitor and return final handler.
func InitMonitor(ctx context.Context, pprofPort int, tracingConfig base.TracingConfig) func() {
	var shutdowns = make(chan func(), 2)

	// Start pprof server if pprofPort is greater than 0.
	if pprofPort > 0 {
		shutdown := startStatsView(pprofPort)
		shutdowns <- shutdown
	}

	// Initialize jaeger tracer if tracing address is set.
	if tracingConfig.Protocol != "" && tracingConfig.Endpoint != "" {
		shutdown, err := initJaegerTracer(ctx, tracingConfig)
		if err != nil {
			logger.Warnf("init jaeger tracer error: %v", err)
			return func() {}
		}

		shutdowns <- shutdown
	}

	return func() {
		logger.Infof("do %d monitor finalizer", len(shutdowns))
		for {
			select {
			case shutdown := <-shutdowns:
				shutdown()
			default:
				return
			}
		}
	}
}

// startStatsView starts the statsview server on the specified port.
func startStatsView(port int) func() {
	addr := net.JoinHostPort(net.IPv4zero.String(), strconv.Itoa(port))
	viewer.SetConfiguration(viewer.WithAddr(addr))
	sv := statsview.New()

	go func() {
		if err := sv.Start(); err != nil {
			logger.Errorf("started statsview on http://%s/debug/statsview error: %v", addr, err)
		}
	}()

	logger.Infof("started statsview on http://%s/debug/statsview", addr)
	return func() {
		logger.Info("stopped statsview")
		sv.Stop()
	}
}

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initJaegerTracer(ctx context.Context, tracingConfig base.TracingConfig) (func(), error) {
	var (
		exporter *otlptrace.Exporter
		err      error
	)

	switch tracingConfig.Protocol {
	case "http":
		exporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(tracingConfig.Endpoint), otlptracehttp.WithURLPath(tracingConfig.Path), otlptracehttp.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("could not create HTTP trace exporter: %w", err)
		}
	case "https":
		exporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(tracingConfig.Endpoint), otlptracehttp.WithURLPath(tracingConfig.Path))
		if err != nil {
			return nil, fmt.Errorf("could not create HTTP trace exporter: %w", err)
		}
	case "grpc":
		conn, err := grpc.NewClient(tracingConfig.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("could not create gRPC connection to collector: %w", err)
		}

		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn), otlptracegrpc.WithHeaders(tracingConfig.Headers))
		if err != nil {
			return nil, fmt.Errorf("could not create gRPC trace exporter: %w", err)
		}
	default:
		panic(fmt.Sprintf("unsupported tracing protocol: %s", tracingConfig.Protocol))
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.HostNameKey.String(fqdn.FQDNHostname),
			semconv.HostIPKey.String(ip.IPv4.String()),
			semconv.ServiceNameKey.String(tracingConfig.ServiceName),
			semconv.ServiceNamespaceKey.String("dragonfly"),
			semconv.ServiceVersionKey.String(version.GitVersion))),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		logger.Info("stoped jaeger tracer")
		if err := provider.Shutdown(ctx); err != nil {
			logger.Errorf("shutdown jaeger tracer error: %v", err)
		}
	}, nil
}

func SetupQuitSignalHandler(handler func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		var done bool
		for {
			select {
			case sig := <-signals:
				logger.Warnf("receive signal: %v", sig)
				if !done {
					done = true
					handler()
					logger.Warnf("handle signal: %v finish", sig)
				}
			}
		}
	}()
}

// initConfig reads in config file and ENV variables if set.
func initConfig(useConfigFile bool, name string, config any) {
	// Use config file and read once.
	if useConfigFile {
		cfgFile := viper.GetString("config")
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(dfpath.DefaultConfigDir)
			viper.SetConfigName(name)
			viper.SetConfigType("yaml")
		}

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			ignoreErr := false
			if errors.As(err, &viper.ConfigFileNotFoundError{}) {
				if cfgFile == "" {
					ignoreErr = true
				}
			}
			if !ignoreErr {
				panic(fmt.Errorf("viper read config: %w", err))
			}
		}
	}
	if err := viper.Unmarshal(config, initDecoderConfig); err != nil {
		panic(fmt.Errorf("unmarshal config to struct: %w", err))
	}
}

func LoadConfig(config any) error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(config, initDecoderConfig)
}

func WatchConfig(interval time.Duration, newConfig func() (cfg any), watcher func(cfg any)) {
	var oldData string
	file := viper.ConfigFileUsed()

	data, err := os.ReadFile(file)
	if err != nil {
		logger.Errorf("read file %s error: %v", file, err)
	}
	oldData = string(data)
loop:
	for {
		select {
		case <-time.After(interval):
			// for k8s configmap case, the config file is symbol link
			// reload file instead use fsnotify
			data, err = os.ReadFile(file)
			if err != nil {
				logger.Errorf("read file %s error: %v", file, err)
				continue loop
			}
			if oldData != string(data) {
				cfg := newConfig()
				err = LoadConfig(cfg)
				if err != nil {
					logger.Errorf("load config file %s error: %v", file, err)
					continue loop
				}
				logger.Infof("config file %s changed", file)
				watcher(cfg)
				oldData = string(data)
			}
		}
	}
}

func initDecoderConfig(dc *mapstructure.DecoderConfig) {
	dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(func(from, to reflect.Type, v any) (any, error) {
		switch to {
		case reflect.TypeOf(dfnet.NetAddr{}),
			reflect.TypeOf(types.PEMContent("")),
			reflect.TypeOf(net.IP{}):

			b, _ := yaml.Marshal(v)
			p := reflect.New(to)
			if err := yaml.Unmarshal(b, p.Interface()); err != nil {
				return nil, err
			}

			return p.Interface(), nil
		default:
			return v, nil
		}
	}, dc.DecodeHook)
}
