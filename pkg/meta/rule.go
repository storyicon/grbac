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
    "github.com/hashicorp/go-multierror"
    jsoniter "github.com/json-iterator/go"
)

// Rules is the list of Rule
type Rules []*Rule

// Rule is used to define the relationship between "resource" and "permission"
type Rule struct {
    // The ID controls the priority of the rule.
    // The higher the ID means the higher the priority of the rule.
    // When a request is matched to more than one rule,
    // then authentication will only use the permission configuration for the rule with the highest ID value.
    // If there are multiple rules that are the largest ID, then one of them will be used randomly.
    ID int `json:"id" yaml:"id"`
    *Resource `yaml:",inline"`
    *Permission `yaml:",inline"`
}

// IsValid is used to test the validity of the Rule
func (rule *Rule) IsValid() error {
    if rule.Resource == nil || rule.Permission == nil {
        return ErrEmptyStructure
    }
    err := rule.Resource.IsValid()
    if err != nil {
        return err
    }
    return rule.Permission.IsValid()
}

// IsValid is used to test the validity of the Rule
func (rules Rules) IsValid() error {
    var errs error
    for _, rule := range rules {
        err := rule.IsValid()
        if err != nil {
            errs = multierror.Append(errs, err)
        }
    }
    if errs != nil {
        return errs
    }
    return nil
}

// IsRolesGranted is used to determine whether the current role is admitted by the current rule.
func (rules Rules) IsRolesGranted(roles []string) (PermissionState, error) {
    if len(rules) == 0 {
        return PermissionNeglected, nil
    }
    tail := rules[0]
    for i := 0; i < len(rules); i++ {
        if tail.ID <= rules[i].ID {
            tail = rules[i]
        }
    }
    return tail.IsGranted(roles)
}

func (rules Rules) String() string {
    s, _ := jsoniter.MarshalToString(rules)
    return s
}
