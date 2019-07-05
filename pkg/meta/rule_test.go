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

import "testing"

func TestRule_IsValid(t *testing.T) {
    type fields struct {
        ID         int
        Resource   *Resource
        Permission *Permission
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        {
            name: "test0",
            fields: fields{
                ID:         0,
                Resource:   nil,
                Permission: nil,
            },
            wantErr: true,
        },
        {
            name: "test1",
            fields: fields{
                ID:         0,
                Resource:   nil,
                Permission: &Permission{},
            },
            wantErr: true,
        },
        {
            name: "test2",
            fields: fields{
                ID:         0,
                Resource:   &Resource{},
                Permission: nil,
            },
            wantErr: true,
        },
        {
            name: "test3",
            fields: fields{
                ID:         0,
                Resource:   &Resource{},
                Permission: &Permission{},
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        rule := &Rule{
            ID:         tt.fields.ID,
            Resource:   tt.fields.Resource,
            Permission: tt.fields.Permission,
        }
        if err := rule.IsValid(); (err != nil) != tt.wantErr {
            t.Errorf("%q. Rule.IsValid() error = %v, wantErr %v", tt.name, err, tt.wantErr)
        }
    }
}

func TestRules_IsValid(t *testing.T) {
    tests := []struct {
        name    string
        rules   Rules
        wantErr bool
    }{
        {
            name:    "test0",
            rules:   Rules{},
            wantErr: false,
        },
        {
            name: "test1",
            rules: Rules{
                {},
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        if err := tt.rules.IsValid(); (err != nil) != tt.wantErr {
            t.Errorf("%q. Rules.IsValid() error = %v, wantErr %v", tt.name, err, tt.wantErr)
        }
    }
}

func TestRules_IsRolesGranted(t *testing.T) {
    type args struct {
        roles []string
    }
    tests := []struct {
        name    string
        rules   Rules
        args    args
        want    PermissionState
        wantErr bool
    }{
        {
            name:    "test0",
            rules:   Rules{},
            args:    args{},
            want:    PermissionNeglected,
            wantErr: false,
        },
        {
            name: "test1",
            rules: Rules{
                {
                    Permission: &Permission{
                        AuthorizedRoles: []string{"visitor"},
                    },
                    Resource: &Resource{Host: "test"},
                },
            },
            args: args{
                roles: []string{"editor"},
            },
            want:    PermissionUngranted,
            wantErr: false,
        },
        {
            name: "test1",
            rules: Rules{
                {
                    Permission: &Permission{
                        AuthorizedRoles: []string{"visitor"},
                    },
                    Resource: &Resource{Host: "test"},
                },
            },
            args: args{
                roles: []string{"editor", "visitor"},
            },
            want:    PermissionGranted,
            wantErr: false,
        },
        {
            name: "test2",
            rules: Rules{
                {
                    Permission: &Permission{
                        AuthorizedRoles: []string{"*"},
                    },
                    Resource: &Resource{Host: "test"},
                },
            },
            args: args{
                roles: []string{"editor", "visitor"},
            },
            want:    PermissionGranted,
            wantErr: false,
        },
        {
            name: "test3",
            rules: Rules{
                {
                    Permission: &Permission{
                        AuthorizedRoles: []string{"*"},
                    },
                    Resource: &Resource{Host: "test"},
                },
            },
            args: args{
                roles: []string{},
            },
            want:    PermissionUngranted,
            wantErr: false,
        },
    }
    for _, tt := range tests {
        got, err := tt.rules.IsRolesGranted(tt.args.roles)
        if (err != nil) != tt.wantErr {
            t.Errorf("%q. Rules.IsRolesGranted() error = %v, wantErr %v", tt.name, err, tt.wantErr)
            continue
        }
        if got != tt.want {
            t.Errorf("%q. Rules.IsRolesGranted() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestRules_String(t *testing.T) {
    tests := []struct {
        name  string
        rules Rules
        want  string
    }{
        {
            name:  "test0",
            rules: Rules{},
            want:  "[]",
        },
    }
    for _, tt := range tests {
        if got := tt.rules.String(); got != tt.want {
            t.Errorf("%q. Rules.String() = %v, want %v", tt.name, got, tt.want)
        }
    }
}
