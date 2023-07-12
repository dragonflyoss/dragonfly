/*
 *     Copyright 2023 The Dragonfly Authors
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

package service

import (
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	trainerv1 "d7y.io/api/pkg/apis/trainer/v1"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/trainer/config"
	"d7y.io/dragonfly/v2/trainer/storage"
)

// V1 is the interface for v1 version of the service.
type V1 struct {
	// Trainer service config.
	config *config.Config

	// Storage Interface.
	storage storage.Storage
}

// New v1 version of service instance.
func NewV1(
	cfg *config.Config,
	storage storage.Storage,

) *V1 {
	return &V1{
		config:  cfg,
		storage: storage,
	}
}

// TODO Implement Train methods of v1 version.
// Train implements the Trainer.Train method.
func (v *V1) Train(stream trainerv1.Trainer_TrainServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&emptypb.Empty{})
			}

			logger.Errorf("receive failed: %s", err.Error())
			return err
		}

		logger := logger.WithTrain(req.Hostname, req.Ip, req.ClusterId)
		switch trainRequest := req.GetRequest().(type) {
		case *trainerv1.TrainRequest_TrainMlpRequest:
		default:
			msg := fmt.Sprintf("receive unknown request: %#v", trainRequest)
			logger.Error(msg)
			return status.Error(codes.FailedPrecondition, msg)
		}
	}
}
