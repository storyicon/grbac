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

func TestResource_Match(t *testing.T) {
    type fields struct {
        Host   string
        Path   string
        Method string
    }
    type args struct {
        query *Query
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    bool
        wantErr bool
    }{
        {
            name: "test0",
            fields: fields{
                Host:   "*",
                Path:   "*",
                Method: "*",
            },
            args: args{
                query: &Query{
                    Host:   "host",
                    Path:   "path",
                    Method: "method",
                },
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "test1",
            fields: fields{
                Host:   "host",
                Path:   "path",
                Method: "method",
            },
            args: args{
                query: &Query{
                    Host:   "host",
                    Path:   "path",
                    Method: "method",
                },
            },
            want:    true,
            wantErr: false,
        },
    }
    for _, tt := range tests {
        r := &Resource{
            Host:   tt.fields.Host,
            Path:   tt.fields.Path,
            Method: tt.fields.Method,
        }
        got, err := r.Match(tt.args.query)
        if (err != nil) != tt.wantErr {
            t.Errorf("%q. Resource.Match() error = %v, wantErr %v", tt.name, err, tt.wantErr)
            continue
        }
        if got != tt.want {
            t.Errorf("%q. Resource.Match() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestResource_GetArguments(t *testing.T) {
    type fields struct {
        Host   string
        Path   string
        Method string
    }
    tests := []struct {
        name   string
        fields fields
        want   []string
    }{
        {
            name: "test0",
            fields: fields{
                Host:   "host",
                Path:   "path",
                Method: "method",
            },
            want: []string{"host", "path", "method"},
        },
        {
            name: "test1",
            fields: fields{
                Host:   "",
                Path:   "",
                Method: "",
            },
            want: []string{"", "", ""},
        },
    }
    for _, tt := range tests {
        r := &Resource{
            Host:   tt.fields.Host,
            Path:   tt.fields.Path,
            Method: tt.fields.Method,
        }
        if got := r.GetArguments(); !reflect.DeepEqual(got, tt.want) {
            t.Errorf("%q. Resource.GetArguments() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestResource_IsValid(t *testing.T) {
    type fields struct {
        Host   string
        Path   string
        Method string
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        {
            name:    "test0",
            fields:  fields{},
            wantErr: true,
        },
        {
            name: "test1",
            fields: fields{
                Host:   "host",
                Path:   "path",
                Method: "method",
            },
            wantErr: false,
        },
        {
            name: "test2",
            fields: fields{
                Host:   "host",
                Path:   "",
                Method: "method",
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        r := &Resource{
            Host:   tt.fields.Host,
            Path:   tt.fields.Path,
            Method: tt.fields.Method,
        }
        if err := r.IsValid(); (err != nil) != tt.wantErr {
            t.Errorf("%q. Resource.IsValid() error = %v, wantErr %v", tt.name, err, tt.wantErr)
        }
    }
}
