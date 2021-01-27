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
	"errors"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc/cdnsystem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pieceSeedStream struct {
	sc   *seederClient
	ctx  context.Context
	sr   *cdnsystem.SeedRequest
	opts []grpc.CallOption

	// client for one target
	client  cdnsystem.SeederClient
	nextNum int
	// stream for one client
	stream cdnsystem.Seeder_ObtainSeedsClient

	rpc.RetryMeta
}

func newPieceSeedStream(sc *seederClient, ctx context.Context, sr *cdnsystem.SeedRequest, opts []grpc.CallOption) (*pieceSeedStream, error) {
	pss := &pieceSeedStream{
		sc:   sc,
		ctx:  ctx,
		sr:   sr,
		opts: opts,

		RetryMeta: rpc.RetryMeta{
			MaxAttempts: 5,
			InitBackoff: 0.5,
			MaxBackOff:  4.0,
		},
	}

	xc, _, nextNum := sc.GetClientSafely()
	pss.client, pss.nextNum = xc.(cdnsystem.SeederClient), nextNum

	if err := pss.initStream(); err != nil {
		return nil, err
	} else {
		return pss, nil
	}
}

func (pss *pieceSeedStream) initStream() error {
	stream, err := rpc.ExecuteWithRetry(func() (interface{}, error) {
		return pss.client.ObtainSeeds(pss.ctx, pss.sr, pss.opts...)
	}, pss.InitBackoff, pss.MaxBackOff, pss.MaxAttempts, nil)

	if err != nil {
		err = pss.replaceClient(err)
	} else {
		pss.stream = stream.(cdnsystem.Seeder_ObtainSeedsClient)
		pss.StreamTimes = 1
	}

	return err
}

func (pss *pieceSeedStream) recv() (ps *cdnsystem.PieceSeed, err error) {
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

	if cause = pss.replaceStream(cause); cause != nil {
		if cause = pss.replaceClient(cause); cause != nil {
			return nil, cause
		}
	}

	return pss.recv()
}

func (pss *pieceSeedStream) replaceStream(cause error) error {
	if pss.StreamTimes >= pss.MaxAttempts {
		return errors.New("times of replacing stream reaches limit")
	}

	stream, cause := rpc.ExecuteWithRetry(func() (interface{}, error) {
		return pss.client.ObtainSeeds(pss.ctx, pss.sr, pss.opts...)
	}, pss.InitBackoff, pss.MaxBackOff, pss.MaxAttempts, cause)

	if cause == nil {
		pss.stream = stream.(cdnsystem.Seeder_ObtainSeedsClient)
		pss.StreamTimes++
	}

	return cause
}

func (pss *pieceSeedStream) replaceClient(cause error) error {
	if cause = pss.sc.TryMigrate(pss.nextNum, cause); cause != nil {
		return cause
	}

	xc, _, nextNum := pss.sc.GetClientSafely()
	pss.client, pss.nextNum = xc.(cdnsystem.SeederClient), nextNum

	stream, cause := rpc.ExecuteWithRetry(func() (interface{}, error) {
		return pss.client.ObtainSeeds(pss.ctx, pss.sr, pss.opts...)
	}, pss.InitBackoff, pss.MaxBackOff, pss.MaxAttempts, cause)

	if cause != nil {
		cause = pss.replaceClient(cause)
	} else {
		pss.stream = stream.(cdnsystem.Seeder_ObtainSeedsClient)
		pss.StreamTimes = 1
	}

	return cause
}
