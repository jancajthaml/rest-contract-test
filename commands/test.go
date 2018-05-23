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
)

func CmdTest(c *cli.Context) error {
	resource := c.Args().First()
	if len(resource) == 0 {
		return fmt.Errorf("no resource provided")
	}

	contract := new(parser.Contract)

	// FIXME determine if resource is url or local file
	err := contract.FromFile(resource)
	if err != nil {
		return err
	}

	fmt.Printf("+------------------------------------------------------------------------\n")
	fmt.Printf("| %s (%s)\n", contract.Source, contract.Type)
	fmt.Printf("+------------------------------------------------------------------------\n")
	fmt.Printf("| title: %s\n", contract.Name)
	fmt.Printf("+------------------------------------------------------------------------\n")
	for _, endpoint := range contract.Endpoints {
		fmt.Printf("| %s | %s\n", endpoint.Method+"      "[0:6-len(endpoint.Method)], endpoint.Path)
	}
	if len(contract.Endpoints) > 0 {
		fmt.Printf("+------------------------------------------------------------------------\n")
	}

	return nil
}
