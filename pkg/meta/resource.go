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
    "github.com/storyicon/grbac/pkg/path"
)

// Resource defines resources
type Resource struct {
    // Host defines the host of the resource, allowing wildcards to be used.
    Host string `json:"host" yaml:"host"`
    // Path defines the path of the resource, allowing wildcards to be used.
    Path string `json:"path" yaml:"path"`
    // Method defines the method of the resource, allowing wildcards to be used.
    Method string `json:"method" yaml:"method"`
}

// Match is used to calculate whether the query matches the resource
func (r *Resource) Match(query *Query) (bool, error) {
    args := query.GetArguments()
    for i, res := range r.GetArguments() {
        matched, err := path.Match(res, args[i])
        if err != nil {
            return false, err
        }
        if !matched {
            return false, nil
        }
    }
    return true, nil
}

// GetArguments is used to convert the current argument to a string slice
func (r *Resource) GetArguments() []string {
    return []string{
        r.Host,
        r.Path,
        r.Method,
    }
}

// IsValid is used to test the validity of the Rule
func (r *Resource) IsValid() error {
    if r.Host == "" || r.Method == "" || r.Path == "" {
        return ErrFieldIncomplete
    }
    return nil
}
