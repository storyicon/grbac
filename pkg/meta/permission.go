// Copyright 2018 storyicon@foxmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package meta

import (
    "github.com/hashicorp/go-multierror"
)

// Permissions is the set of Permission
type Permissions []*Permission

// Permission is used to define permission control information
type Permission struct {
    // AuthorizedRoles defines roles that allow access to specified resource
    // Accepted type: non-empty string, *
    //      *: means any role, but visitors should have at least one role,
    //      non-empty string: specified role
    AuthorizedRoles []string `json:"authorized_roles" yaml:"authorized_roles"`
    // ForbiddenRoles defines roles that not allow access to specified resource
    // ForbiddenRoles has a higher priority than AuthorizedRoles
    // Accepted type: non-empty string, *
    //      *: means any role, but visitors should have at least one role,
    //      non-empty string: specified role
    //
    ForbiddenRoles []string `json:"forbidden_roles" yaml:"forbidden_roles"`
    // AllowAnyone has a higher priority than ForbiddenRoles/AuthorizedRoles
    // If set to true, anyone will be able to pass authentication.
    // Note that this will include people without any role.
    AllowAnyone bool `json:"allow_anyone" yaml:"allow_anyone"`
}

// IsValid is used to test the validity of the Rule
func (p *Permission) IsValid() error {
    if p.AllowAnyone == false && len(p.AuthorizedRoles) == 0 && len(p.ForbiddenRoles) == 0 {
        return multierror.Prefix(ErrEmptyStructure, "permission: ")
    }
    return nil
}

// IsGranted is used to determine whether the given role can pass the authentication of *Permission.
func (p *Permission) IsGranted(roles []string) (PermissionState, error) {
    if p.AllowAnyone {
        return PermissionGranted, nil
    }

    if len(roles) == 0 {
        return PermissionUngranted, nil
    }

    for _, role := range roles {
        for _, forbidden := range p.ForbiddenRoles {
            if forbidden == "*" || (role == forbidden) {
                return PermissionUngranted, nil
            }
        }
        for _, authorized := range p.AuthorizedRoles {
            if authorized == "*" || (role == authorized) {
                return PermissionGranted, nil
            }
        }
    }
    return PermissionUngranted, nil
}
