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

package path

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestHasWildcardPrefix(t *testing.T) {
    type args struct {
        pattern string
    }
    tests := []struct {
        name string
        args args
        want bool
    }{
    // TODO: Add test cases.
    }
    for _, tt := range tests {
        if got := HasWildcardPrefix(tt.args.pattern); got != tt.want {
            t.Errorf("%q. HasWildcardPrefix() = %v, want %v", tt.name, got, tt.want)
        }
    }
}

func TestTrimWildcard(t *testing.T) {
    type args struct {
        pattern string
    }
    tests := []struct {
        name            string
        args            args
        wantTrimmed     string
        wantHasWildcard bool
    }{
    // TODO: Add test cases.
    }
    for _, tt := range tests {
        gotTrimmed, gotHasWildcard := TrimWildcard(tt.args.pattern)
        if gotTrimmed != tt.wantTrimmed {
            t.Errorf("%q. TrimWildcard() gotTrimmed = %v, want %v", tt.name, gotTrimmed, tt.wantTrimmed)
        }
        if gotHasWildcard != tt.wantHasWildcard {
            t.Errorf("%q. TrimWildcard() gotHasWildcard = %v, want %v", tt.name, gotHasWildcard, tt.wantHasWildcard)
        }
    }
}

func TestMatch(t *testing.T) {
    type Result struct {
        Matched bool
        Err bool
    }
    var TestMatchEqual = func(wanted Result, pattern, s string) {
        matched, err := Match(pattern, s)
        result := Result{matched, err != nil}
        assert.Equal(t, wanted, result, fmt.Sprintf("[Match(..) != Wanted] Match(..):%+v, Wanted:%+v, pattern: %s, s: %s ", result, wanted, pattern, s))
    }
    TestMatchEqual(Result{true,  false}, `*`, ``)
    TestMatchEqual(Result{false,  false}, `*`, `/`) // Wrong
    TestMatchEqual(Result{false,  false}, `/*`, `//`) // Wrong
    TestMatchEqual(Result{true,  false}, `*/`, `debug/`)
    TestMatchEqual(Result{true,  false}, `/*`, `/debug`)
    TestMatchEqual(Result{false,  false}, `/*`, `/debug/`)
    TestMatchEqual(Result{false,  false}, `/*`, `/debug/`) // Wrong
    TestMatchEqual(Result{false,  false}, `/*`, `/debug/pprof`)
    TestMatchEqual(Result{true,  false}, `/*/`, `/debug/`)
    TestMatchEqual(Result{true,  false}, `/*/*`, `/debug/pprof`)
    TestMatchEqual(Result{true,  false}, `debug/*/`, `debug/test/`)
    TestMatchEqual(Result{true  ,  false}, `aa/*`, `aa/`) // Wrong
    TestMatchEqual(Result{true,  false}, `**`, ``)
    TestMatchEqual(Result{true,  false}, `/**`, `/debug`)
    TestMatchEqual(Result{true,  false}, `/**`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true,  false}, `/**`, `/debug/pprof/profile/`)
    TestMatchEqual(Result{true,  false}, `/in[d]ex`, `/index`)
    TestMatchEqual(Result{false,  false}, `/in[d]ex`, `/inex`)
    TestMatchEqual(Result{true,  false}, `/in\[d\]ex`, `/in[d]ex`)
    TestMatchEqual(Result{false,  false}, `/**/profile`, `/debug/pprof/profile/`) // Wrong
    TestMatchEqual(Result{true,  false}, `/**/profile`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true,  false}, `/*/*/profile`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true,  false}, `/**/*`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true,  false}, `/**/pprof/*`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true,  false}, `/**/pprof/*/`, `/debug/pprof/profile/`)
    TestMatchEqual(Result{true,  false}, `/*/[pz]rofile/`, `/debug/profile/`)
    TestMatchEqual(Result{true,  false}, `/{debug,test}/profile`, `/debug/profile`)
    TestMatchEqual(Result{false,  false}, `/{debug,test}/profile`, `/debug/profile/`)
    TestMatchEqual(Result{true,  false}, `\**`, `*GET`)
    TestMatchEqual(Result{true,  false}, `\\[0-9]`, `\8`)
    TestMatchEqual(Result{true,  false}, `\\\[0-9]`, `\[0-9]`)
    TestMatchEqual(Result{true,  false}, `\\A`, `\A`)
    TestMatchEqual(Result{true,  false}, `\A`, `A`)
    TestMatchEqual(Result{false,  false}, `[^visitor]*`, `va`)
    TestMatchEqual(Result{true,  false}, `dashboard*.xxxx.com`, `dashboard.xxxx.com`)
    TestMatchEqual(Result{true,  false}, `dashboard{-sit,-prod}.xxxx.com`, `dashboard-sit.xxxx.com`)
    TestMatchEqual(Result{false,  false}, `dashboard{-sit,-prod}.xxxx.com`, `dashboard-si.xxxx.com`)
    TestMatchEqual(Result{false,  true}, `/{config/*,instance}`, `/config/delete`)
    TestMatchEqual(Result{true,  false}, `**`, `/config`)
}

