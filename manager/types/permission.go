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

package types

type RoleParams struct {
	RoleName string `uri:"role_name" binding:"required,min=1"`
}

type UserParams struct {
	UserName string `uri:"user_name" binding:"required"`
}

type ObjectPermission struct {
	Object string `json:"object" binding:"required,min=1"`
	Action string `json:"aciton" binding:"omitempty,oneof=read *"`
}

type RolePermissionCreateRequest struct {
	RoleName   string             `json:"role_name" binding:"required"`
	Permissons []ObjectPermission `json:"permissions" binding:"dive"`
}

type RolePermissionUpdateRequest struct {
	Method     string             `json:"method" binding:"required,oneof=add remove"`
	Permissons []ObjectPermission `json:"permissions" binding:"dive"`
}

type Permissions []string
