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

func TestQuery_GetArguments(t *testing.T) {
    tests := []struct {
        name  string
        query *Query
        want  []string
    }{
        {
            name: "test0",
            query: &Query{
                Host:   "host",
                Path:   "path",
                Method: "method",
            },
            want: []string{"host", "path", "method"},
        },
    }
    for _, tt := range tests {
        if got := tt.query.GetArguments(); !reflect.DeepEqual(got, tt.want) {
            t.Errorf("%q. Query.GetArguments() = %v, want %v", tt.name, got, tt.want)
        }
    }
}
