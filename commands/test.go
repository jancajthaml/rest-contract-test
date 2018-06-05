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

	"github.com/jancajthaml/rest-contract-test/parser"
	"github.com/jancajthaml/rest-contract-test/model"

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
	
	//fmt.Printf("Source: %s\n", contract.Source)
	//fmt.Printf("Type: %s\n", contract.Type)

	for _, endpoint := range contract.Endpoints {
		for _, curl := range GenerateCurls(endpoint) {
			fmt.Println(curl)
		}
	}

	return nil
}

// FIXME for debugging right now

func GenerateCurls(ref model.Endpoint) []string {
	res := make([]string, 0)

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

	for k, v := range ref.Headers {
		cmd += "-H \"" + k + ": " + v + "\" "
	}

	if len(ref.Requests) == 0 {
		return append(res, cmd + ref.URI + qs)
	}

	for mime, payload := range ref.Requests {
		// FOR form data: -d "param1=value1&param2=value2"
		switch mime {
		case "application/json":
			bytes, err := json.Marshal(payload)
		    if err != nil {
		        fmt.Println("Can't serialize", err)
		        continue
		    }
			res = append(res, cmd + "-H \"Content-Type: " + mime + "\" " + "-H \"Accept: application/json\" " + "-d '" + string(bytes) + "' " + ref.URI + qs)		
		}
	}

	return res
}
