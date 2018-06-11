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

type Payload struct {
	Content *Content
	Headers map[string]string
}

type Content struct {
	Type    string
	Example interface{}
	Schema  interface{}
}

type Endpoint struct {
	URI          string
	Method       string
	Responses    map[int]Payload
	Request      Payload
	QueryStrings map[string]string
	Provides     Set
	Requires     Set
}

type Contract struct {
	Source    string
	Type      string
	Name      string
	Endpoints []*Endpoint
}
