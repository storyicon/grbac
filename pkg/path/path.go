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
    "strings"

    "github.com/storyicon/grbac/pkg/path/doublestar"
)

// HasWildcardPrefix is​used to determine whether an expression is a wildcard at the beginning
func HasWildcardPrefix(pattern string) bool {
    if len(pattern) == 0 {
        return false
    }
    switch pattern[0] {
    case '?', '*', '[', '{':
        return true
    }
    return false
}

// TrimWildcard is used to intercept the pattern before the first wildcard
func TrimWildcard(pattern string) (trimmed string, hasWildcard bool) {
    var chars []byte
Pattern:
    for i := 0; i < len(pattern); i++ {
        switch pattern[i] {
        case '\\':
            if i == len(pattern)-1 {
                break Pattern
            }
            i++
        case '?', '*', '[', '{':
            hasWildcard = true
            break Pattern
        }
        chars = append(chars, pattern[i])
    }
    return string(chars), hasWildcard
}

// Match returns true if name matches the shell file name pattern.
// The pattern syntax is:
//
//  pattern:
//    { term }
//  term:
//    '*'         matches any sequence of non-path-separators
//    '**'        matches any sequence of characters, including
//                path separators.
//    '?'         matches any single non-path-separator character
//    '[' [ '^' ] { character-range } ']'
//          character class (must be non-empty)
//    '{' { term } [ ',' { term } ... ] '}'
//    c           matches character c (c != '*', '?', '\\', '[')
//    '\\' c      matches character c
//
//  character-range:
//    c           matches character c (c != '\\', '-', ']')
//    '\\' c      matches character c
//    lo '-' hi   matches character c for lo <= c <= hi
//
// Match requires pattern to match all of name, not just a substring.
// The path-separator defaults to the '/' character. The only possible
// returned error is ErrBadPattern, when pattern is malformed.
//
// Note: this is meant as a drop-in replacement for path.Match() which
// always uses '/' as the path separator. If you want to support systems
// which use a different path separator (such as Windows), what you want
// is the PathMatch() function below.
//
func Match(pattern string, s string) (bool, error) {
    switch pattern {
    case "**":
        return true, nil
    case "*":
        if strings.Contains(s, "/") {
            return false, nil
        }
        return true, nil
    }
    return doublestar.Match(pattern, s)
}
