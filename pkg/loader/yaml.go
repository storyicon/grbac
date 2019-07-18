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
    "io/ioutil"

    "github.com/storyicon/grbac/pkg/meta"
    "gopkg.in/yaml.v3"
)

// YAMLLoader implements the Loader interface
// it is used to load configuration from a local yaml file.
type YAMLLoader struct {
    path string
}

// NewYAMLLoader is used to initialize a YAMLLoader
func NewYAMLLoader(file string) (*YAMLLoader, error) {
    loader := &YAMLLoader{
        path: file,
    }
    _, err := loader.Load()
    if err != nil {
        return nil, err
    }
    return loader, nil
}

// Load is used to return a list of rules
func (loader *YAMLLoader) Load() (meta.Rules, error) {
    bytes, err := ioutil.ReadFile(loader.path)
    if err != nil {
        return nil, err
    }
    rules := meta.Rules{}
    err = yaml.Unmarshal(bytes, &rules)
    if err != nil {
        return nil, err
    }
    return rules, nil
}
