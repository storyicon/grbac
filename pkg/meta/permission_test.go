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
    "reflect"
    "testing"
)

func TestPermission_IsValid(t *testing.T) {
    tests := []struct {
        name    string
        p       *Permission
        wantErr bool
    }{
        {
            name: "test0",
            p: &Permission{
                AuthorizedRoles: []string{},
                ForbiddenRoles:  []string{},
                AllowAnyone:     false,
            },
            wantErr: true,
        },
        {
            name: "test1",
            p: &Permission{
                AuthorizedRoles: []string{},
                ForbiddenRoles:  []string{},
                AllowAnyone:     true,
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        if err := tt.p.IsValid(); (err != nil) != tt.wantErr {
            t.Errorf("%q. Permission.IsValid() error = %v, wantErr %v", tt.name, err, tt.wantErr)
        }
    }
}

func TestPermission_IsGranted(t *testing.T) {
    type args struct {
        roles []string
    }
    tests := []struct {
        name    string
        p       *Permission
        args    args
        want    PermissionState
        wantErr bool
    }{
        {
            name: "test0",
            p: &Permission{
                AllowAnyone: true,
            },
            args:    args{},
            want:    PermissionGranted,
            wantErr: false,
        },
        {
            name: "test1",
            p: &Permission{
                AuthorizedRoles: []string{"editor"},
                AllowAnyone:     false,
            },
            args:    args{roles: []string{"editor"}},
            want:    PermissionGranted,
            wantErr: false,
        },
    }
    for _, tt := range tests {
        got, err := tt.p.IsGranted(tt.args.roles)
        if (err != nil) != tt.wantErr {
            t.Errorf("%q. Permission.IsGranted() error = %v, wantErr %v", tt.name, err, tt.wantErr)
            continue
        }
        if !reflect.DeepEqual(got, tt.want) {
            t.Errorf("%q. Permission.IsGranted() = %v, want %v", tt.name, got, tt.want)
        }
    }
}
