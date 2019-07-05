/*
 * Copyright 2019 storyicon@foxmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grbac

import (
    "github.com/storyicon/grbac/pkg/meta"
)

// Resource defines resources
type Resource = meta.Resource

// PermissionState identifies the status of the permission
type PermissionState = meta.PermissionState

// Permissions is the set of Permission
type Permissions = meta.Permissions

// Permission is used to define permission control information
type Permission = meta.Permission

// Rules is the list of Rule
type Rules = meta.Rules

// Rule is used to define the relationship between "resource" and "permission"
type Rule = meta.Rule

// Query defines the data structure of the query parameters
type Query = meta.Query
