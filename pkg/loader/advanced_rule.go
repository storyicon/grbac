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

package loader

import (
    "github.com/storyicon/grbac/pkg/meta"
)

// AdvancedRule allows you to write RBAC rules in a more concise way
type AdvancedRule struct {
    Host   []string `json:"host"`
    Path   []string `json:"path"`
    Method []string `json:"method"`

    *meta.Permission
}

// AdvancedRules is the list of AdvancedRules
type AdvancedRules []*AdvancedRule

// GetRules is used to convert AdvancedRules to meta.Rules
func (adv AdvancedRules) GetRules() meta.Rules {
    var rules meta.Rules
    for _, item := range adv {
        for _, host := range item.Host {
            for _, path := range item.Path {
                for _, method := range item.Method {
                    rules = append(rules, &meta.Rule{
                        Resource: &meta.Resource{
                            Host:   host,
                            Path:   path,
                            Method: method,
                        },
                        Permission: item.Permission,
                    })
                }
            }
        }
    }
    return rules
}

// AdvancedRulesLoader implements the Loader interface
// it is used to load configuration from advanced data.
type AdvancedRulesLoader struct {
    rules AdvancedRules
}

// NewAdvancedRulesLoader is used to initialize a AdvancedRulesLoader
func NewAdvancedRulesLoader(rules AdvancedRules) (*AdvancedRulesLoader, error) {
    return &AdvancedRulesLoader{
        rules: rules,
    }, nil
}

// Load is used to return a list of rules
func (loader *AdvancedRulesLoader) Load() (meta.Rules, error) {
    return loader.rules.GetRules(), nil
}
