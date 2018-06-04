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

type Response struct {
	Example string // FIXME example is a wrong name
	Schema  string // FIXME not string but interface{} instead ?
}

type Request struct {
	Example string // FIXME example is a wrong name
	Schema  string // FIXME not string but interface{} instead ?
}

type Endpoint struct {
	URI          string
	Method       string
	Responses    []Response // FIXME map HTTP_CODE -> RESPONSE
	Headers      map[string]string
	Request      Request // FIXME optional ?
	QueryStrings map[string]string
}

type Contract struct {
	Source    string
	Type      string
	Name      string
	Endpoints []Endpoint
}

func (ref Endpoint) String() string {
	qs := Urlencode(ref.QueryStrings)
	if len(qs) != 0 {
		qs = "?" + qs
	}

	cmd := "curl -v -L "

	switch ref.Method {
	case "PUT":
		cmd += "-X PUT "
	case "POST":
		cmd += "-X POST "
	case "PATCH":
		cmd += "-X PATCH "
	case "DELETE":
		cmd += "-X DELETE "
	}

	for k, v := range ref.Headers {
		cmd += "-H \"" + k + ": " + v + "\" "
	}

	//--data '{"username":"xyz","password":"xyz"}' \

	return cmd + ref.URI + qs
}
