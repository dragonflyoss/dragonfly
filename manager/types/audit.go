/*
 *     Copyright 2025 The Dragonfly Authors
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

package types

import "time"

type CreateAuditRequest struct {
	ActorType  string    `json:"actor_type" binding:"required"`
	ActorName  string    `json:"actor_name" binding:"required"`
	EventType  string    `json:"event_type" binding:"required"`
	Operation  string    `json:"operation" binding:"required"`
	OperatedAt time.Time `json:"operated_at" binding:"required"`
	State      string    `json:"status" binding:"required"`
	Path       string    `json:"path" binding:"omitempty"`
	StatusCode int       `json:"status_code" binding:"omitempty"`
}

type GetAuditsQuery struct {
	ActorType  string `form:"actor_type" binding:"omitempty,oneof=UNKNOWN USER PAT"`
	ActorName  string `form:"actor_name" binding:"omitempty"`
	EventType  string `form:"event_type" binding:"omitempty,oneof=API"`
	Operation  string `form:"operation" binding:"omitempty"`
	State      string `form:"state" binding:"omitempty,oneof=SUCCESS FAILURE"`
	Path       string `form:"path" binding:"omitempty"`
	StatusCode int    `form:"status_code" binding:"omitempty"`
	Page       int    `form:"page" binding:"omitempty,gte=1"`
	PerPage    int    `form:"per_page" binding:"omitempty,gte=1,lte=10000000"`
}
