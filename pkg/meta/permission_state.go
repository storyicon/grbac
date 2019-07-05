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

// PermissionState identifies the status of the permission
type PermissionState uint8

const (
    // PermissionUnknown is an initial state, usually specified when an error occurs
    PermissionUnknown PermissionState = iota
    // PermissionGranted means permission is granted
    PermissionGranted
    // PermissionUngranted means permission is ungranted
    PermissionUngranted
    // PermissionNeglected means could not find the matching rule in the list of rules
    PermissionNeglected
)

// IsLooselyGranted is used to determine whether a request is authorized in a non-strict sense
// It returns true when state equals PermissionGranted or PermissionNeglected
// * This means if you forget to configure some addresses, they may be accessed by anyone.
func (state PermissionState) IsLooselyGranted() bool {
    return (state == PermissionGranted) || (state == PermissionNeglected)
}

// IsNeglected is used to determine if the current state is equal to PermissionNeglected
// PermissionNeglected means could not find the matching rule in the list of rules
func (state PermissionState) IsNeglected() bool {
    return state == PermissionNeglected
}

// IsGranted is used to determine whether the current request is granted in a strict sense.
// Note that it only returns true when state equals PermissionGranted
// Because we recommend that you configure permissions for all possible requests to prevent forgetting to configure some addresses
// * If you want it to return true when PermissionNeglected as well, you should use IsLooselyGranted
func (state PermissionState) IsGranted() bool {
    return state == PermissionGranted
}

func (state PermissionState) String() string {
    switch state {
    case PermissionGranted:
        return "Permission Granted"
    case PermissionUngranted:
        return "Permission Ungranted"
    case PermissionNeglected:
        return "Permission Neglected"
    default:
        return "Permission Unknown"
    }
}
