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

    jsoniter "github.com/json-iterator/go"
    "github.com/storyicon/grbac/pkg/meta"
)

// JSONLoader implements the Loader interface
// it is used to load configuration from a local json file.
type JSONLoader struct {
    path string
}

// NewJSONLoader is used to initialize a JSONLoader
func NewJSONLoader(file string) (*JSONLoader, error) {
    loader := &JSONLoader{
        path: file,
    }
    _, err := loader.Load()
    if err != nil {
        return nil, err
    }
    return loader, nil
}

// Load is used to return a list of rules
func (loader *JSONLoader) Load() (meta.Rules, error) {
    bytes, err := ioutil.ReadFile(loader.path)
    if err != nil {
        return nil, err
    }
    rules := meta.Rules{}
    err = jsoniter.Unmarshal(bytes, &rules)
    if err != nil {
        return nil, err
    }
    return rules, nil
}
