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

package server

import (
	"net/http"

	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/database"
	"d7y.io/dragonfly/v2/manager/service"

	// manager server rpc
	_ "d7y.io/dragonfly/v2/pkg/rpc/manager/server"
)

type Server struct {
	config     *config.Config
	service    *service.Service
	restServer *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	// ms, err := service.NewManagerServer(cfg)
	// if err != nil {
	// return nil, fmt.Errorf("Failed to create manager server: %s", err)
	// }

	// router, err := initRouter(ms)
	// if err != nil {
	// return nil, err
	// }

	// return &Server{
	// config:     cfg,
	// rpcService: ms,
	// restServer: &http.Server{
	// Addr:    ":8080",
	// Handler: router,
	// },
	// }, nil

	db, err := database.New(cfg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Server) Serve() error {
	// go func() {
	// port := s.cfg.Server.Port
	// err := rpc.StartTCPServer(port, port, s.rpcService)
	// if err != nil {
	// logger.Errorf("failed to start manager tcp server: %+v", err)
	// }

	// s.stop <- struct{}{}
	// }()

	// go func() {
	// if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// logger.Errorf("failed to start manager http server: %+v", err)
	// }

	// s.stop <- struct{}{}
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// select {
	// case <-s.stop:
	// s.Stop()
	// case <-quit:
	// s.Stop()
	// }

	return nil
}

func (s *Server) Stop() {
	// if s.ms != nil {
	// err := s.ms.Close()
	// if err != nil {
	// logger.Errorf("failed to stop manager server: %+v", err)
	// }

	// s.ms = nil
	// }

	// rpc.StopServer()

	// ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	// defer cancel()

	// err := s.httpServer.Shutdown(ctx)
	// if err != nil {
	// logger.Errorf("failed to stop manager http server: %+v", err)
	// }
}
