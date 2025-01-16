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

package main

import (
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
)

type evaluator struct{}

// EvaluateParents sort parents by evaluating multiple feature scores.
func (e *evaluator) EvaluateParents(parents []*standard.Peer, child *standard.Peer, taskPieceCount int32) []*standard.Peer {
	return []*standard.Peer{&standard.Peer{}}
}

// IsBadParent determine if peer is a bad parent, it can not be selected as a parent.
func (e *evaluator) IsBadParent(peer *standard.Peer) bool {
	return true
}

// EvaluatePersistentCacheParents sort persistent cache parents by evaluating multiple feature scores.
func (e *evaluator) EvaluatePersistentCacheParents(parents []*persistentcache.Peer, child *persistentcache.Peer, taskPieceCount int32) []*persistentcache.Peer {
	return []*persistentcache.Peer{&persistentcache.Peer{}}
}

// IsBadParent determine if persistent cache peer is a bad parent, it can not be selected as a parent.
func (e *evaluator) IsBadPersistentCacheParent(peer *persistentcache.Peer) bool {
	return true
}

func DragonflyPluginInit(option map[string]string) (any, map[string]string, error) {
	return &evaluator{}, map[string]string{"type": "scheduler", "name": "evaluator"}, nil
}
