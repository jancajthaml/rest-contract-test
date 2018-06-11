// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package model

import "fmt"

// Set is set datastructure representing distinct slice
type Set struct {
	items map[string]interface{}
}

// NewSet returns empty set
func NewSet() Set {
	return Set{make(map[string]interface{})}
}

// Add adds element to set
func (set *Set) Add(i string) {
	set.items[i] = nil
}

// AddAll adds all elements of set into this set
func (set *Set) AddAll(input Set) {
	for k := range input.items {
		set.items[k] = nil
	}
}

// Contains returns true if value is present in set
func (set *Set) Contains(i string) bool {
	_, found := set.items[i]
	return found
}

// Remove removes element from set
func (set *Set) Remove(i string) {
	delete(set.items, i)
}

// Size returns number of items in set
func (set *Set) Size() int {
	return len(set.items)
}

// AsSlice returns set as slice
func (set *Set) AsSlice() []string {
	keys := make([]string, 0, len(set.items))
	for k := range set.items {
		keys = append(keys, k)
	}
	return keys
}

// Copy returns deep copy
func (set *Set) Copy() Set {
	clone := make(map[string]interface{})
	for k := range set.items {
		clone[k] = nil
	}
	return Set{items: clone}
}

// String returns string representation of set
func (set Set) String() string {
	return fmt.Sprintf("%v", set.AsSlice())
}
