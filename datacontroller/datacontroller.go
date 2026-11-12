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

// Package datacontroller runs the Kubernetes controllers that manage Dataset and
// DataLifecyclePolicy objects on top of a Dragonfly manager.
package datacontroller

import (
	"context"
	"fmt"

	"github.com/go-logr/zapr"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
	"d7y.io/dragonfly/v2/datacontroller/config"
	"d7y.io/dragonfly/v2/datacontroller/controllers"
	logger "d7y.io/dragonfly/v2/internal/dflog"
)

// Server is the data controller process.
type Server struct {
	config  *config.Config
	manager ctrl.Manager
}

// New builds the controller manager and registers the reconcilers.
func New(cfg *config.Config) (*Server, error) {
	ctrl.SetLogger(zapr.NewLogger(logger.CoreLogger.Desugar()))

	restConfig, err := restConfigFor(cfg)
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	resync := cfg.Controller.ResyncInterval
	options := ctrl.Options{
		Scheme:                  scheme,
		HealthProbeBindAddress:  cfg.Controller.HealthProbeAddr,
		LeaderElection:          cfg.Controller.LeaderElection,
		LeaderElectionID:        cfg.Controller.LeaderElectionID,
		LeaderElectionNamespace: cfg.Controller.LeaderElectionNamespace,
		Cache:                   cache.Options{SyncPeriod: &resync},
		Metrics:                 metricsserver.Options{BindAddress: "0"},
	}

	if cfg.Metrics.Enable {
		options.Metrics.BindAddress = cfg.Metrics.Addr
	}

	if cfg.Controller.Namespace != "" {
		options.Cache.DefaultNamespaces = map[string]cache.Config{cfg.Controller.Namespace: {}}
	}

	mgr, err := ctrl.NewManager(restConfig, options)
	if err != nil {
		return nil, fmt.Errorf("create controller manager: %w", err)
	}

	datasetReconciler := &controllers.DatasetReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorder("dataset-controller"),
		Config:   cfg,
	}
	if err := datasetReconciler.SetupWithManager(mgr); err != nil {
		return nil, fmt.Errorf("setup dataset controller: %w", err)
	}

	policyReconciler := &controllers.DataLifecyclePolicyReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	if err := policyReconciler.SetupWithManager(mgr); err != nil {
		return nil, fmt.Errorf("setup data lifecycle policy controller: %w", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("add healthz check: %w", err)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("add readyz check: %w", err)
	}

	return &Server{config: cfg, manager: mgr}, nil
}

// Serve runs the controllers until the context is cancelled.
func (s *Server) Serve(ctx context.Context) error {
	logger.Infof("starting data controller, manager endpoint %q, namespace %q, leader election %t",
		s.config.Manager.Endpoint, s.config.Controller.Namespace, s.config.Controller.LeaderElection)
	return s.manager.Start(ctx)
}

// restConfigFor loads the kubeconfig named in the configuration, falling back to the
// standard lookup (KUBECONFIG, in-cluster).
func restConfigFor(cfg *config.Config) (*rest.Config, error) {
	if cfg.Controller.KubeConfig != "" {
		restConfig, err := clientcmd.BuildConfigFromFlags("", cfg.Controller.KubeConfig)
		if err != nil {
			return nil, fmt.Errorf("load kubeconfig %q: %w", cfg.Controller.KubeConfig, err)
		}

		return restConfig, nil
	}

	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("load kubernetes configuration: %w", err)
	}

	return restConfig, nil
}
