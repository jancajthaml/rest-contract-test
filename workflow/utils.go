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

func discoverRequestRequirements(request interface{}, requirements *[]string) {
	switch x := request.(type) {

	case map[string]interface{}:
		for _, v := range x {
			discoverRequestRequirements(v, requirements)
		}

	case string:
		for _, submatches := range placeholderPattern.FindAllStringSubmatch(x, -1) {
			for _, match := range submatches {
				*requirements = append(*requirements, match)
			}
		}

	case map[interface{}]interface{}:
		for _, v := range x {
			discoverRequestRequirements(v, requirements)
		}

	case []interface{}:
		for _, v := range x {
			discoverRequestRequirements(v, requirements)
		}

	case []string:
		for _, v := range x {
			discoverRequestRequirements(v, requirements)
		}

	}
	return
}

func PopulateRequirements(contract *model.Contract) {
	for _, endpoint := range contract.Endpoints {

		// uri requirements
		for _, submatches := range placeholderPattern.FindAllStringSubmatch(endpoint.URI, -1) {
			for _, match := range submatches {
				endpoint.Requires = append(endpoint.Requires, match)
			}
		}

		// headers requirements
		for _, val := range endpoint.Headers {
			for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
				for _, match := range submatches {
					endpoint.Requires = append(endpoint.Requires, match)
				}
			}
		}

		// queryString requirements
		for _, val := range endpoint.QueryStrings {
			for _, submatches := range placeholderPattern.FindAllStringSubmatch(val, -1) {
				for _, match := range submatches {
					endpoint.Requires = append(endpoint.Requires, match)
				}
			}
		}

		// request requirements
		discoverRequestRequirements(endpoint.Request, &endpoint.Requires)

		if len(endpoint.Requires) != 0 {
			fmt.Println("endpoint", endpoint.Method, endpoint.URI, "requires following placeholders:", endpoint.Requires)
		}
	}

	return
}

func PopulateProvisions(contract *model.Contract) {
	// FIXME separate ?

	fmt.Println("populating provisions")

	globals := make([]string, 0)

	for _, pair := range os.Environ() {
		providing := strings.Split(pair, "=")[0]
		alias := strings.Replace(strings.ToLower(providing), "_", "-", -1)
		globals = append(globals, providing)
		if alias != providing {
			globals = append(globals, alias)
		}
	}

	//if len(globals) != 0 {
	//fmt.Println("globals provide:", globals)
	//}

	// TBD then response provisions ... info based on response code (prefix?)

	return
}
