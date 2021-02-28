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
	"d7y.io/dragonfly/v2/pkg/basic/dfnet"
	"d7y.io/dragonfly/v2/pkg/rpc/base"
	"d7y.io/dragonfly/v2/pkg/rpc/base/common"
	cdnclient "d7y.io/dragonfly/v2/pkg/rpc/cdnsystem/client"
	"d7y.io/dragonfly/v2/pkg/rpc/scheduler"
	"fmt"
	"google.golang.org/grpc"
	"strings"
)


func GetPieceTasks(destPeer *scheduler.PeerPacket_DestPeer, ctx context.Context, ptr *base.PieceTaskRequest, opts ...grpc.CallOption) (*base.PiecePacket, error) {
	destAddr := fmt.Sprintf("%s:%d", destPeer.Ip, destPeer.RpcPort)
	peerId := destPeer.PeerId
	toCdn := strings.HasSuffix(peerId, common.CdnSuffix)
	var err error
	client, err := getClient(destAddr, toCdn)
	if err != nil {
		return nil, err
	}
	if toCdn {
		return client.(cdnclient.SeederClient).GetPieceTasks(ctx, ptr, opts...)
	}
	return client.(DaemonClient).GetPieceTasks(ctx, ptr, opts...)
}

func getClient(destAddr string, toCdn bool) (interface{}, error) {
	netAddr := dfnet.NetAddr{
		Type: dfnet.TCP,
		Addr: destAddr,
	}

	if toCdn {
		return cdnclient.GetClientByAddr([]dfnet.NetAddr{netAddr})
	} else {
		return GetClientByAddr([]dfnet.NetAddr{netAddr})
	}
}
