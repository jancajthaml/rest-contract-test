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

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	gio "github.com/jancajthaml/rest-contract-test/io"
)

var placeholderPattern = regexp.MustCompile(`(?:\{|\<{2}).{1,100}?(?:\}|\>{2})`)

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

// FIXME add MarkFailed and MarkSuccessfull methods to endpoint, add *bool variable
// representing if it was successfull or failed (or <nil> which means skipped)

type Contract struct {
	Source    string
	Type      string
	Name      string
	Endpoints []*Endpoint
}

func (ref *Endpoint) Prepare(variables map[string]string) error {

	//fmt.Println(">>> variables", variables.AsSlice())

	//fmt.Println("preparing endpoint")
	//fmt.Println("before", ref)

	// uri requirements
	for _, submatches := range placeholderPattern.FindAllStringSubmatch(ref.URI, -1) {
		for _, match := range submatches {
			if rv, ok := variables[match]; ok {
				fmt.Println("satisfied", match, "in uri", rv)
				ref.URI = strings.Replace(ref.URI, match, rv, -1)
			} else {
				return fmt.Errorf("unsatisfied requirement %s", match)
			}
		}

	}

	// queryString requirements
	for k, val := range ref.QueryStrings {
		for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
			for _, match := range submatches {
				if rv, ok := variables[match]; ok {
					fmt.Println("satisfied", match, "in queryString", rv)
					ref.QueryStrings[k] = strings.Replace(val, match, rv, -1)
				} else {
					return fmt.Errorf("unsatisfied requirement %s", match)
				}
			}
		}
	}

	// headers requirements
	for k, val := range ref.Request.Headers {
		for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
			for _, match := range submatches {
				if rv, ok := variables[match]; ok {
					fmt.Println("satisfied", match, "in headers", rv)
					ref.Request.Headers[k] = strings.Replace(val, match, rv, -1)
				} else {
					return fmt.Errorf("unsatisfied requirement %s", match)
				}
			}
		}
	}

	return nil
}

func (ref Endpoint) Mark(err error) {
	if err != nil {
		fmt.Println("Failed", ref.Method, ref.URI, "with", err)
	} else {
		fmt.Println("Success", ref.Method, ref.URI)
	}
}

func (ref Endpoint) React(variables map[string]string, code int, respContent []byte) {

	switch code {

	case 404, 405:
		ref.Mark(fmt.Errorf("Invalid call %d", code))

	default:
		response, ok := ref.Responses[code]
		if !ok {
			ref.Mark(fmt.Errorf("Undocumented response code %d", code))
			return
		}

		if response.Content != nil {
			switch response.Content.Type {
			case "application/json":
				var body interface{}
				if err := json.Unmarshal(respContent, &body); err != nil {
					ref.Mark(err)
					return
				}

				body = gio.UntypedConvert(body)

				fmt.Println("body of response is", body, "with reference of", response.Content.Example)
			}
		}

		ref.Mark(nil)

	}

	return
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

	for k, v := range ref.Request.Headers {
		cmd += "-H \"" + k + ": " + v + "\" "
	}

	if ref.Request.Content != nil {
		switch ref.Request.Content.Type {
		case "application/json":
			if bytes, err := json.Marshal(ref.Request.Content.Example); err == nil {
				cmd += "-H \"Content-Type: " + ref.Request.Content.Type + "\" "
				cmd += "-H \"Accept: application/json\" "
				cmd += "-d '" + string(bytes) + "' "
			}
		}
	}

	return cmd + ref.URI + qs
}
