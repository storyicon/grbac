// Copyright (c) 2019 https://github.com/bmatcuk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package path

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
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
        {
            name: "test0",
            args: args{
                pattern: "*",
            },
            want: true,
        },
        {
            name: "test1",
            args: args{
                pattern: "jack*",
            },
            want: false,
        },
        {
            name: "test2",
            args: args{
                pattern: `\*tom`,
            },
            want: false,
        },
        {
            name: "test3",
            args: args{
                pattern: "/test",
            },
            want: false,
        },
        {
            name: "test4",
            args: args{
                pattern: "[t]est",
            },
            want: true,
        },
        {
            name: "test5",
            args: args{
                pattern: "{t,j}est",
            },
            want: true,
        },
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
        {
            name: "test0",
            args: args{
                pattern: "*test",
            },
            wantTrimmed:     "",
            wantHasWildcard: true,
        },
        {
            name: "test1",
            args: args{
                pattern: "test*",
            },
            wantTrimmed:     "test",
            wantHasWildcard: true,
        },
        {
            name: "test2",
            args: args{
                pattern: "te*st",
            },
            wantTrimmed:     "te",
            wantHasWildcard: true,
        },
        {
            name: "test3",
            args: args{
                pattern: "test",
            },
            wantTrimmed:     "test",
            wantHasWildcard: false,
        },
        {
            name: "test4",
            args: args{
                pattern: `test\[]`,
            },
            wantTrimmed:     `test[]`,
            wantHasWildcard: false,
        },
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
        Err     bool
    }
    var TestMatchEqual = func(wanted Result, pattern, s string) {
        matched, err := Match(pattern, s)
        result := Result{matched, err != nil}
        assert.Equal(t, wanted, result, fmt.Sprintf("[Match(..) != Wanted] Match(..):%+v, Wanted:%+v, pattern: %s, s: %s ", result, wanted, pattern, s))
    }
    TestMatchEqual(Result{true, false}, `*`, ``)
    TestMatchEqual(Result{false, false}, `*`, `/`)   // Wrong
    TestMatchEqual(Result{false, false}, `/*`, `//`) // Wrong
    TestMatchEqual(Result{true, false}, `*/`, `debug/`)
    TestMatchEqual(Result{true, false}, `/*`, `/debug`)
    TestMatchEqual(Result{false, false}, `/*`, `/debug/`)
    TestMatchEqual(Result{false, false}, `/*`, `/debug/`) // Wrong
    TestMatchEqual(Result{false, false}, `/*`, `/debug/pprof`)
    TestMatchEqual(Result{true, false}, `/*/`, `/debug/`)
    TestMatchEqual(Result{true, false}, `/*/*`, `/debug/pprof`)
    TestMatchEqual(Result{true, false}, `debug/*/`, `debug/test/`)
    TestMatchEqual(Result{true, false}, `aa/*`, `aa/`) // Wrong
    TestMatchEqual(Result{true, false}, `**`, ``)
    TestMatchEqual(Result{true, false}, `/**`, `/debug`)
    TestMatchEqual(Result{true, false}, `/**`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true, false}, `/**`, `/debug/pprof/profile/`)
    TestMatchEqual(Result{true, false}, `/in[d]ex`, `/index`)
    TestMatchEqual(Result{false, false}, `/in[d]ex`, `/inex`)
    TestMatchEqual(Result{true, false}, `/in\[d\]ex`, `/in[d]ex`)
    TestMatchEqual(Result{false, false}, `/**/profile`, `/debug/pprof/profile/`) // Wrong
    TestMatchEqual(Result{true, false}, `/**/profile`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true, false}, `/*/*/profile`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true, false}, `/**/*`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true, false}, `/**/pprof/*`, `/debug/pprof/profile`)
    TestMatchEqual(Result{true, false}, `/**/pprof/*/`, `/debug/pprof/profile/`)
    TestMatchEqual(Result{true, false}, `/*/[pz]rofile/`, `/debug/profile/`)
    TestMatchEqual(Result{true, false}, `/{debug,test}/profile`, `/debug/profile`)
    TestMatchEqual(Result{false, false}, `/{debug,test}/profile`, `/debug/profile/`)
    TestMatchEqual(Result{true, false}, `\**`, `*GET`)
    TestMatchEqual(Result{true, false}, `\\[0-9]`, `\8`)
    TestMatchEqual(Result{true, false}, `\\\[0-9]`, `\[0-9]`)
    TestMatchEqual(Result{true, false}, `\\A`, `\A`)
    TestMatchEqual(Result{true, false}, `\A`, `A`)
    TestMatchEqual(Result{false, false}, `[^visitor]*`, `va`)
    TestMatchEqual(Result{true, false}, `dashboard*.xxxx.com`, `dashboard.xxxx.com`)
    TestMatchEqual(Result{true, false}, `dashboard{-sit,-prod}.xxxx.com`, `dashboard-sit.xxxx.com`)
    TestMatchEqual(Result{false, false}, `dashboard{-sit,-prod}.xxxx.com`, `dashboard-si.xxxx.com`)
    TestMatchEqual(Result{false, true}, `/{config/*,instance}`, `/config/delete`)
    TestMatchEqual(Result{true, false}, `**`, `/config`)
}
