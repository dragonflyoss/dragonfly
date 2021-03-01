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

package client

import (
	"context"
	"d7y.io/dragonfly/v2/pkg/rpc"
	"d7y.io/dragonfly/v2/pkg/rpc/cdnsystem"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pieceSeedStream struct {
	// client for one target
	sc      *seederClient
	ctx     context.Context
	hashKey string
	sr      *cdnsystem.SeedRequest
	opts    []grpc.CallOption
	// stream for one client
	stream cdnsystem.Seeder_ObtainSeedsClient
	// server list which cannot serve
	failedServers []string
	rpc.RetryMeta
}

func newPieceSeedStream(sc *seederClient, ctx context.Context, hashKey string, sr *cdnsystem.SeedRequest, opts []grpc.CallOption) (*pieceSeedStream, error) {
	pss := &pieceSeedStream{
		sc:      sc,
		ctx:     ctx,
		hashKey: hashKey,
		sr:      sr,
		opts:    opts,
		RetryMeta: rpc.RetryMeta{
			MaxAttempts: 5,
			InitBackoff: 0.5,
			MaxBackOff:  4.0,
		},
	}

	if err := pss.initStream(); err != nil {
		return nil, err
	} else {
		return pss, nil
	}
}

func (pss *pieceSeedStream) initStream() error {
	stream, err := rpc.ExecuteWithRetry(func() (interface{}, error) {
		if client, err := pss.sc.getSeederClient(pss.hashKey); err != nil {
			return nil, err
		} else {
			return client.ObtainSeeds(pss.ctx, pss.sr, pss.opts...)
		}
	}, pss.InitBackoff, pss.MaxBackOff, pss.MaxAttempts, nil)
	if err != nil {
		err = pss.replaceClient(pss.hashKey, err)
	} else {
		pss.stream = stream.(cdnsystem.Seeder_ObtainSeedsClient)
		pss.StreamTimes = 1
	}
	return err
}

func (pss *pieceSeedStream) recv() (ps *cdnsystem.PieceSeed, err error) {
	pss.sc.UpdateAccessNodeMap(pss.hashKey)
	if ps, err = pss.stream.Recv(); err != nil {
		ps, err = pss.retryRecv(err)
	}
	return
}

func (pss *pieceSeedStream) retryRecv(cause error) (*cdnsystem.PieceSeed, error) {
	code := status.Code(cause)
	if code == codes.DeadlineExceeded {
		return nil, cause
	}

	if err := pss.replaceStream(pss.hashKey, cause); err != nil {
		if err := pss.replaceClient(pss.hashKey, cause); err != nil {
			return nil, cause
		}
	}

	return pss.recv()
}

func (pss *pieceSeedStream) replaceStream(key string, cause error) error {
	if pss.StreamTimes >= pss.MaxAttempts {
		return errors.New("times of replacing stream reaches limit")
	}

	stream, err := rpc.ExecuteWithRetry(func() (interface{}, error) {
		if client, err := pss.sc.getSeederClient(key); err != nil {
			return nil, err
		} else {
			return client.ObtainSeeds(pss.ctx, pss.sr, pss.opts...)
		}
	}, pss.InitBackoff, pss.MaxBackOff, pss.MaxAttempts, cause)

	if err == nil {
		pss.stream = stream.(cdnsystem.Seeder_ObtainSeedsClient)
		pss.StreamTimes++
	}

	return err
}

func (pss *pieceSeedStream) replaceClient(key string, cause error) error {
	if preNode, err := pss.sc.TryMigrate(key, cause, pss.failedServers); err != nil {
		return err
	} else {
		pss.failedServers = append(pss.failedServers, preNode)
	}
	stream, err := rpc.ExecuteWithRetry(func() (interface{}, error) {
		if client, err := pss.sc.getSeederClient(key); err != nil {
			return nil, err
		} else {
			return client.ObtainSeeds(pss.ctx, pss.sr, pss.opts...)
		}
	}, pss.InitBackoff, pss.MaxBackOff, pss.MaxAttempts, cause)

	if err != nil {
		err = pss.replaceClient(key, cause)
	} else {
		pss.stream = stream.(cdnsystem.Seeder_ObtainSeedsClient)
		pss.StreamTimes = 1
	}

	return err
}
