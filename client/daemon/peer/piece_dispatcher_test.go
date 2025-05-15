/*
 *     Copyright 2022 The Dragonfly Authors
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

package peer

import (
	"cmp"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
)

type peerDesc struct {
	id         string
	uploadTime time.Duration
	count      *int
}

type pieceTestManager struct {
	pieceDispatcher PieceDispatcher
	peers           map[string]peerDesc
	pieceNum        int
}

func newPieceTestManager(pieceDispatcher PieceDispatcher, peers []peerDesc, pieceNum int) *pieceTestManager {
	peerMap := make(map[string]peerDesc)
	for _, p := range peers {
		peerMap[p.id] = peerDesc{
			id:         p.id,
			uploadTime: p.uploadTime,
			count:      new(int),
		}
	}
	pm := &pieceTestManager{
		pieceDispatcher: pieceDispatcher,
		peers:           peerMap,
		pieceNum:        pieceNum,
	}
	return pm
}

func (pc *pieceTestManager) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// producer
	go func() {
		slice := make([]*DownloadPieceRequest, 0)
		for range 4 {
			for _, peer := range pc.peers {
				for j := range pc.pieceNum {
					slice = append(slice, &DownloadPieceRequest{
						piece:  &commonv1.PieceInfo{PieceNum: int32(j)},
						DstPid: peer.id,
					})
				}
			}
		}
		rand.Shuffle(len(slice), func(i, j int) {
			tmp := slice[i]
			slice[i] = slice[j]
			slice[j] = tmp
		})
		for _, req := range slice {
			pc.pieceDispatcher.Put(req)
		}
		wg.Done()
	}()

	// consumer
	wg.Add(1)
	go func() {
		downloaded := map[int32]struct{}{}
		for {
			req, err := pc.pieceDispatcher.Get()
			if err == ErrNoValidPieceTemporarily {
				continue
			}
			if err != nil {
				break
			}
			*pc.peers[req.DstPid].count = *pc.peers[req.DstPid].count + 1
			downloaded[req.piece.PieceNum] = struct{}{}
			pc.pieceDispatcher.Report(&DownloadPieceResult{
				pieceInfo:  req.piece,
				DstPeerID:  req.DstPid,
				BeginTime:  time.Now().UnixNano(),
				FinishTime: time.Now().Add(pc.peers[req.DstPid].uploadTime).UnixNano(),
			})
			if len(downloaded) >= pc.pieceNum {
				break
			}
		}
		pc.pieceDispatcher.Close()
		wg.Done()
	}()
	wg.Wait()
}

// return a slice of peerIDs, order by count desc
func (pc *pieceTestManager) Order() []string {
	peerIDs := maps.Keys(pc.peers)
	slices.SortFunc(peerIDs, func(a, b string) int { return cmp.Compare(*pc.peers[a].count, *pc.peers[b].count) })
	return peerIDs
}
