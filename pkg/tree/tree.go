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

// Tree is a Radix search tree that supports wildcards
type Tree struct {
    root *Node
}

// NewTree is used to initialize a wildcard tree
func NewTree() *Tree {
    root := NewNode("ROOT", nil)
    return &Tree{
        root: root,
    }
}

// Query is used to query the current tree by args
func (tree *Tree) Query(args []string) ([]Data, error) {
    var data []Data

    parents := []*Node{tree.root}
    for i, arg := range args {
        eof := i == len(args)-1

        var nodes []*Node
        for _, parent := range parents {
            children, childData, err := parent.Find(arg)
            if err != nil {
                return nil, err
            }
            nodes = append(nodes, children...)
            if eof {
                data = append(data, childData...)
            }
        }
        parents = nodes
    }
    return data, nil
}

// Insert is used to insert a node into the current tree
func (tree *Tree) Insert(args []string, data Data) {
    parent := tree.root

    var nodeData Data
    for i, arg := range args {
        eof := i == len(args)-1
        if eof {
            nodeData = data
        }
        child := NewNode(arg, nodeData)
        parent.Insert(child)
        parent = child
    }
}
