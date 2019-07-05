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

func TestPermissionState_IsLooselyGranted(t *testing.T) {
    tests := []struct {
        name  string
        state PermissionState
        want  bool
    }{
        {
            name:  "test0",
            state: PermissionNeglected,
            want:  true,
        },
    }
    for _, tt := range tests {
        if got := tt.state.IsLooselyGranted(); got != tt.want {
            t.Errorf("%q. PermissionState.IsLooselyGranted() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestPermissionState_IsNeglected(t *testing.T) {
    tests := []struct {
        name  string
        state PermissionState
        want  bool
    }{
        {
            name:  "test0",
            state: PermissionUngranted,
            want:  false,
        },
        {
            name:  "test1",
            state: PermissionNeglected,
            want:  true,
        },
    }
    for _, tt := range tests {
        if got := tt.state.IsNeglected(); got != tt.want {
            t.Errorf("%q. PermissionState.IsNeglected() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestPermissionState_IsGranted(t *testing.T) {
    tests := []struct {
        name  string
        state PermissionState
        want  bool
    }{
        {
            name:  "test0",
            state: PermissionUngranted,
            want:  false,
        },
        {
            name:  "test1",
            state: PermissionUnknown,
            want:  false,
        },
        {
            name:  "test2",
            state: PermissionNeglected,
            want:  false,
        },
    }
    for _, tt := range tests {
        if got := tt.state.IsGranted(); got != tt.want {
            t.Errorf("%q. PermissionState.IsGranted() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestPermissionState_String(t *testing.T) {
    tests := []struct {
        name  string
        state PermissionState
        want  string
    }{
        {
            name:  "test0",
            state: PermissionNeglected,
            want:  "Permission Neglected",
        },
    }
    for _, tt := range tests {
        if got := tt.state.String(); got != tt.want {
            t.Errorf("%q. PermissionState.String() = %v, want %v", tt.name, got, tt.want)
        }
    }
}
