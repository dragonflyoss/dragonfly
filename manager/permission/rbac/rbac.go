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

package rbac

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	managermodels "d7y.io/dragonfly/v2/manager/models"
)

// Syntax for models see https://casbin.org/docs/en/syntax-for-models
const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && (r.act == p.act || p.act == "*")
`

const (
	RootRole  = "root"
	GuestRole = "guest"
)

const (
	AllAction  = "*"
	ReadAction = "read"
)

const (
	// RootUserName is the name of the built-in root user.
	RootUserName = "root"

	// DragonflyInitialRootPasswordEnvName is the environment variable name of the root
	// user's initial password. It is only used when seeding the root user for the first
	// time; the password is managed in the database afterwards.
	DragonflyInitialRootPasswordEnvName = "DRAGONFLY_INITIAL_ROOT_PASSWORD"

	// DefaultRootPassword is the initial password of the root user when
	// DragonflyInitialRootPasswordEnvName is not set.
	DefaultRootPassword = "dragonfly"
)

// MinRootPasswordLength and MaxRootPasswordLength match the password binding of
// types.SignInRequest, so that the seeded password can always be used to sign in.
const (
	MinRootPasswordLength = 8
	MaxRootPasswordLength = 20
)

var (
	apiGroupRegexp = regexp.MustCompile(`^/api/v[0-9]+/([-_a-zA-Z]*)[/.*]*`)
)

// initialRootPassword returns the password to seed the root user with, which is
// DragonflyInitialRootPasswordEnvName if set and DefaultRootPassword otherwise.
func initialRootPassword() (string, error) {
	password := os.Getenv(DragonflyInitialRootPasswordEnvName)
	if password == "" {
		return DefaultRootPassword, nil
	}

	if len(password) < MinRootPasswordLength || len(password) > MaxRootPasswordLength {
		return "", fmt.Errorf("%s must be between %d and %d characters", DragonflyInitialRootPasswordEnvName, MinRootPasswordLength, MaxRootPasswordLength)
	}

	return password, nil
}

func NewEnforcer(gdb *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(gdb, &managermodels.CasbinRule{})
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

func InitRBAC(e *casbin.Enforcer, g *gin.Engine, db *gorm.DB) error {
	// Create roles
	permissions := GetPermissions(g)
	for _, permission := range permissions {
		if _, err := e.AddPermissionForUser(RootRole, permission.Object, AllAction); err != nil {
			return err
		}

		if _, err := e.AddPermissionForUser(GuestRole, permission.Object, ReadAction); err != nil {
			return err
		}
	}

	// Create root user for the first time
	var rootUserCount int64
	if err := db.Model(managermodels.User{}).Count(&rootUserCount).Error; err != nil {
		return err
	}

	if rootUserCount <= 0 {
		password, err := initialRootPassword()
		if err != nil {
			return err
		}

		if password == DefaultRootPassword {
			logger.Warnf("seed the root user with the default password, set %s to override it", DragonflyInitialRootPasswordEnvName)
		}

		encryptedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			return err
		}

		rootUser := managermodels.User{
			EncryptedPassword: string(encryptedPasswordBytes),
			Name:              RootUserName,
			State:             managermodels.UserStateEnabled,
		}

		if err := db.Create(&rootUser).Error; err != nil {
			return err
		}

		if _, err := e.AddRoleForUser(fmt.Sprint(rootUser.ID), RootRole); err != nil {
			return err
		}
	}

	return nil
}

type Permission struct {
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required,oneof=read *"`
}

func GetPermissions(g *gin.Engine) []Permission {
	permissions := []Permission{}
	actions := []string{AllAction, ReadAction}
	for _, permission := range GetAPIGroupNames(g) {
		for _, action := range actions {
			permissions = append(permissions, Permission{
				Object: permission,
				Action: action,
			})
		}
	}

	return permissions
}

func GetAPIGroupNames(g *gin.Engine) []string {
	apiGroupNames := []string{}
	for _, route := range g.Routes() {
		name, err := GetAPIGroupName(route.Path)
		if err != nil {
			continue
		}

		if !slices.Contains(apiGroupNames, name) {
			apiGroupNames = append(apiGroupNames, name)
		}
	}

	return apiGroupNames
}

func GetAPIGroupName(path string) (string, error) {
	matchs := apiGroupRegexp.FindStringSubmatch(path)
	if len(matchs) != 2 {
		return "", errors.New("cannot find group name")
	}

	return matchs[1], nil
}

func HTTPMethodToAction(method string) string {
	action := ReadAction
	if method == http.MethodDelete || method == http.MethodPatch || method == http.MethodPut || method == http.MethodPost {
		action = AllAction
	}

	return action
}
