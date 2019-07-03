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

package grbac

import (
    "github.com/storyicon/grbac/pkg/meta"
    "github.com/stretchr/testify/assert"
    "testing"
)

type Result struct {
    State PermissionState
    Error error
}

func NewQuery(c *Controller, host, path, method string, roles []string) *Result {
    state, err := c.IsQueryGranted(&meta.Query{
        Host: host,
        Path: path,
        Method: method,
    }, roles)
    return &Result{
        State: state,
        Error: err,
    }
}

func TestNew(t *testing.T) {
    var rules Rules
    rules = append(rules,
        &Rule{ Resource: &Resource{ Host: `*`, Path: `**`, Method: `*`, }, Permission: &Permission{ AuthorizedRoles: []string{"*"}, ForbiddenRoles: []string{"black_user"}, AllowAnyone: false, }},
        &Rule{ Resource: &Resource{ Host: `domain.com`, Path: `**`, Method: `*`, }, Permission: &Permission{ AuthorizedRoles: []string{}, ForbiddenRoles: []string{}, AllowAnyone: true, }},
        &Rule{ Resource: &Resource{ Host: `dashboard-{prod,sit}.domain.com`, Path: `/config`, Method: `*`, }, Permission: &Permission{ AuthorizedRoles: []string{"sre"}, ForbiddenRoles: []string{}, AllowAnyone: false, }},
        &Rule{ Resource: &Resource{ Host: `pprof.domain.com`, Path: `/**`, Method: `*`, }, Permission: &Permission{ AuthorizedRoles: []string{"engineer", "sre"}, ForbiddenRoles: []string{}, AllowAnyone: false, }},
        &Rule{ Resource: &Resource{ Host: `pprof.domain.com`, Path: `/virtual/*`, Method: `*`, }, Permission: &Permission{ AuthorizedRoles: []string{}, ForbiddenRoles: []string{}, AllowAnyone: true, }},
        &Rule{ Resource: &Resource{ Host: `domain.com`, Path: `/api/**`, Method: `POST`, }, Permission: &Permission{ AuthorizedRoles: []string{"editor", "engineer"}, ForbiddenRoles: []string{}, AllowAnyone: false, }},
        &Rule{ Resource: &Resource{ Host: `domain.com`, Path: `/api/**`, Method: `DELETE`, }, Permission: &Permission{ AuthorizedRoles: []string{"super_editor"}, ForbiddenRoles: []string{}, AllowAnyone: false, }},
        &Rule{ Resource: &Resource{ Host: `x-domain.com`, Path: `/articles`, Method: `{DELETE,POST,PUT}`, }, Permission: &Permission{ AuthorizedRoles: []string{"editor"}, ForbiddenRoles: []string{}, AllowAnyone: false, }},
        &Rule{ Resource: &Resource{ Host: `x-domain.com`, Path: `/articles`, Method: `{GET}`, }, Permission: &Permission{ AuthorizedRoles: []string{}, ForbiddenRoles: []string{}, AllowAnyone: true, }},
    )

    c, err := New(WithRules(rules))
    assert.Equal(t, nil, err)

    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "domain.com", "/index.html", "GET", []string{}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "dashboard-prod.domain.com", "/index.html", "GET", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "dashboard-prod.domain.com", "/index.html", "GET", []string{"black_user"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "dashboard-prod.domain.com", "/index.html", "GET", []string{}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "dashboard-sit.domain.com", "/index.html", "GET", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "pprof.domain.com", "/index.html", "GET", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "pprof.domain.com", "/index.html", "GET", []string{"engineer"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "pprof.domain.com", "/virtual/f91fj2f1rj043rj9043e21esfhasdh09a", "GET", []string{}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "pprof.domain.com", "/config/get", "GET", []string{"sre"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "pprof.domain.com", "/config/get", "GET", []string{"anyone"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "domain.com", "/api/get", "POST", []string{"anyone"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "domain.com", "/api/get", "POST", []string{}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "domain.com", "/api/get", "POST", []string{"editor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "domain.com", "/api/get", "POST", []string{"engineer"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "domain.com", "/api/get", "DELETE", []string{"anyone"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "domain.com", "/api/get", "DELETE", []string{"super_editor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "DELETE", []string{"editor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "POST", []string{"editor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "PUT", []string{"editor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "GET", []string{"editor"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "DELETE", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "POST", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionUngranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "PUT", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "GET", []string{"visitor"}))
    assert.Equal(t, &Result{State: meta.PermissionGranted, Error: nil}, NewQuery(c, "x-domain.com", "/articles", "GET", []string{}))
}
