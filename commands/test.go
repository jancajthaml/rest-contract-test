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

	"github.com/jancajthaml/rest-contract-test/http"
	"github.com/jancajthaml/rest-contract-test/model"
	"github.com/jancajthaml/rest-contract-test/parser"
	"github.com/jancajthaml/rest-contract-test/workflow"
)

func CmdTest(c *cli.Context) error {
	resource := c.Args().First()
	if len(resource) == 0 {
		return fmt.Errorf("no resource provided")
	}

	contract, err := parser.FromResource(resource)
	if err != nil {
		fmt.Println("resource loaded with error", err)
		return err
	}

	client := http.NewHttpClient()

	Sort(contract)

	for _, endpoint := range contract.Endpoints {
		err := client.Call(endpoint)
		if err != nil {
			fmt.Println("ERROR |", *endpoint)
		}
	}

	return nil
}

func Sort(contract *model.Contract) {
	doneRequirements := make(chan bool)
	doneProvisions := make(chan bool)

	go func() {
		workflow.PopulateRequirements(contract)
		doneRequirements <- true
	}()
	go func() {
		workflow.PopulateProvisions(contract)
		doneProvisions <- true
	}()

	<-doneRequirements
	<-doneProvisions

	workflow.SortEndpoints(contract)
}
