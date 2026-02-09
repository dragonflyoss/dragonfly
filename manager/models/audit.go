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

package models

import "time"

const (
	// ActorTypeUnknown represents an unknown actor type.
	// e.g if user not authorized, the actor type will be UNKNOWN.
	ActorTypeUnknown = "UNKNOWN"

	// ActorTypeUser represents a user actor type.
	// e.g if user authorized by JWT token, the actor type will be USER.
	ActorTypeUser = "USER"

	// ActorTypePat represents a personal access token actor type.
	// e.g if user authorized by personal access token, the actor type will be PAT.
	ActorTypePat = "PAT"
)

const (
	// EventTypeAPI represents the event is from API operation.
	EventTypeAPI = "API"
)

const (
	// AuditStateSuccess represents the successful state.
	AuditStateSuccess = "SUCCESS"

	// AuditStateFailure represents the failure state.
	AuditStateFailure = "FAILURE"
)

type Audit struct {
	BaseModel
	// ActorType represents the actor type, which can be one of the following:
	// - UNKNOWN: Represents an unknown actor such as not authenticated.
	// - USER: Represents a user.
	// - PAT: Represents a personal access token.
	ActorType string `gorm:"column:actor_type;type:varchar(256);not null;comment:actor type" json:"actor_type"`
	// ActorName represents the actor name, it can be the username or token name.
	ActorName string `gorm:"column:actor_name;type:varchar(256);not null;comment:actor name" json:"actor_name"`
	// EventType represents the event type, indicates the type of event, API for http request,
	// can expand to other types such as SYSTEM for internal system events in future.
	EventType string `gorm:"column:event_type;type:varchar(100);not null;comment:event type" json:"event_type"`
	// Operation represents the operation, it will be the HTTP method for API events.
	Operation string `gorm:"column:operation;type:varchar(256);not null;comment:operation" json:"operation"`
	// OperatedAt represents the operation time.
	OperatedAt time.Time `gorm:"column:operated_at;type:timestamp;default:current_timestamp;comment:operated at" json:"operated_at"`
	// State represents the state, it indicates the state of the operation, e.g SUCCESS for API status code >= 200 & < 300.
	State string `gorm:"column:state;type:varchar(100);not null;comment:state" json:"state"`
	// Path represents the request path, it will be the URL path for API events.
	Path string `gorm:"column:path;type:varchar(1024);comment:request path" json:"path"`
	// StatusCode represents the status code, can be ignored for non-API events.
	StatusCode int `gorm:"column:status_code;type:int;comment:status code" json:"status_code"`
	// Detail represents the detail, leave for extension for future use.
	Detail JSONMap `gorm:"column:detail;comment:detail" json:"detail"`
}
