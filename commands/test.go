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

package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/jancajthaml/rest-contract-test/model"
	"github.com/jancajthaml/rest-contract-test/parser"
	"github.com/jancajthaml/rest-contract-test/workflow"

	"encoding/json"
)

func CmdTest(c *cli.Context) error {
	resource := c.Args().First()
	if len(resource) == 0 {
		return fmt.Errorf("no resource provided")
	}

	contract, err := parser.FromResource(resource)
	if err != nil {
		return err
	}

	// FIXME in parallel
	workflow.PopulateRequirements(contract)
	workflow.PopulateProvisions(contract)

	// FIXME wait here

	for _, endpoint := range contract.Endpoints {
		fmt.Println(endpoint.Method, endpoint.URI, "requires:", endpoint.Requires, "provides:", endpoint.Provides)
	}

	return nil
}

// FIXME for debugging right now

func GenerateCurl(ref model.Endpoint) string {
	//res := make([]string, 0)

	qs := model.Urlencode(ref.QueryStrings)
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
