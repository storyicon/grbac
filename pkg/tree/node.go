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

package tree

import (
    iradix "github.com/hashicorp/go-immutable-radix"
    "github.com/storyicon/grbac/pkg/path"
)

// Node defines the wildcard node
type Node struct {
    key           string
    indexKey      []byte
    isWildcardKey bool

    data Data

    tree     *iradix.Tree
    catchAll []*Node
}

// Data is the data type of the data node
type Data = interface{}

// NewNode is used to create a new node
func NewNode(key string, data Data) *Node {
    trimmed, isWildcardKey := path.TrimWildcard(key)
    return &Node{
        key:           key,
        indexKey:      []byte(trimmed),
        isWildcardKey: isWildcardKey,
        data:          data,
        tree:          iradix.New(),
        catchAll:      []*Node{},
    }
}

// Match is used to determine whether the current node's key matches the given key.
func (node *Node) match(key string) (bool, error) {
    if node.isWildcardKey {
        return path.Match(node.key, key)
    }
    return node.key == key, nil
}

// Find is used to find child nodes by a specified key
func (node *Node) Find(key string) ([]*Node, []Data, error) {

    nodes := node.catchAll
    node.tree.Root().WalkPath([]byte(key), func(k []byte, v interface{}) bool {
        children, ok := v.([]*Node)
        if ok {
            nodes = append(nodes, children...)
            return false
        }
        return true
    })

    var tmp []*Node
    var data []Data
    for _, node := range nodes {
        matched, err := node.match(key)
        if err != nil {
            return nil, nil, err
        }
        if matched {
            if node.data != nil {
                data = append(data, node.data)
            }
            tmp = append(tmp, node)
        }
    }

    return tmp, data, nil
}

// Insert used to insert a node into the child node of the current node
func (node *Node) Insert(child *Node) {
    if path.HasWildcardPrefix(child.key) {
        node.catchAll = append(node.catchAll, child)
    } else {
        nodeData, exists := node.tree.Get(child.indexKey)
        nodes := []*Node{child}
        if exists {
            children, _ := nodeData.([]*Node)
            nodes = append(children, child)
        }
        node.tree, _, _ = node.tree.Insert(child.indexKey, nodes)
    }
}
