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

package workflow

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jancajthaml/rest-contract-test/model"
)

var placeholderPattern = regexp.MustCompile(`(?:\{|\<{2}).{1,100}?(?:\}|\>{2})`)

func discoverContentPlaceholders(request interface{}, set *model.Set) {
	switch x := request.(type) {

	case map[string]interface{}:
		for _, v := range x {
			discoverContentPlaceholders(v, set)
		}

	case string:
		for _, submatches := range placeholderPattern.FindAllStringSubmatch(x, -1) {
			for _, match := range submatches {
				set.Add(match)
			}
		}

	case map[interface{}]interface{}:
		for _, v := range x {
			discoverContentPlaceholders(v, set)
		}

	case []interface{}:
		for _, v := range x {
			discoverContentPlaceholders(v, set)
		}

	case []string:
		for _, v := range x {
			discoverContentPlaceholders(v, set)
		}

	}
	return
}

func PopulateRequirements(contract *model.Contract) {
	globals := make(map[string]string)

	// environment provisions
	for _, pair := range os.Environ() {
		providing := strings.Split(pair, "=")

		alias := strings.Replace(strings.ToLower(providing[0]), "_", "-", -1)
		globals["<<"+alias+">>"] = providing[1]
		globals["{"+alias+"}"] = providing[1]
		// FIXME maybe alias also to camelCase and PascalCase for old-time folks
	}

	for _, endpoint := range contract.Endpoints {
		endpoint.Requires = model.NewSet()

		// uri requirements
		for _, submatches := range placeholderPattern.FindAllStringSubmatch(endpoint.URI, -1) {
			for _, match := range submatches {
				if rv, ok := globals[match]; ok {
					endpoint.URI = strings.Replace(endpoint.URI, match, rv, -1)
					continue
				}
				endpoint.Requires.Add(match)
			}

		}

		// queryString requirements
		for _, val := range endpoint.QueryStrings {
			for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
				for _, match := range submatches {
					if rv, ok := globals[match]; ok {
						fmt.Println("queryString requirement would be satisfied by env", match)
						fmt.Println("now must replace value", match, "in", val, "with", rv)
					} else {
						endpoint.Requires.Add(match)
					}
					//endpoint.Requires.Add(match)
					//fmt.Println("queryString", match)
				}
			}
		}

		// headers requirements
		for _, val := range endpoint.Request.Headers {
			for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
				for _, match := range submatches {
					//endpoint.Requires.Add(match)
					//fmt.Println("headers", match)

					if rv, ok := globals[match]; ok {
						fmt.Println("headers requirement would be satisfied by env", match)
						fmt.Println("now must replace value", match, "in", val, "with", rv)
					} else {
						endpoint.Requires.Add(match)
					}
				}
			}
		}

		// request requirements
		if endpoint.Request.Content != nil {
			discoverContentPlaceholders(endpoint.Request.Content.Example, &endpoint.Requires)
		}
	}

	return
}

func PopulateProvisions(contract *model.Contract) {
	// responses provisions
	for _, endpoint := range contract.Endpoints {
		endpoint.Provides = model.NewSet()

	inner:
		for code, response := range endpoint.Responses {
			if code != 200 {
				continue inner
			}

			// response headers provisions
			for _, val := range response.Headers {
				for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
					for _, match := range submatches {
						endpoint.Provides.Add(match)
					}
				}
			}

			// response content provisions
			if response.Content != nil {
				discoverContentPlaceholders(response.Content.Example, &endpoint.Provides)
			}
		}

	}

	// TBD then response provisions ... info based on response code (prefix?)

	return
}

func CalculateOrdering(contract *model.Contract) {

	//for _, endpoint := range contract.Endpoints {
	//obtainables.AddAll(endpoint.Provides)
	//}

	/*
		fmt.Println(obtainables.AsSlice())

		ordering_clean := make([]string, 0)
		ordering_volatile := make([]string, 0)

		satisfied_variables := model.NewSet()

	*/

	// FIXME TBD

	return
}
