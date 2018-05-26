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

package parser

import (
	"fmt"

	"github.com/jancajthaml/rest-contract-test/io"
	"github.com/jancajthaml/rest-contract-test/model"
	"github.com/jancajthaml/rest-contract-test/parser/raml/v04"
	"github.com/jancajthaml/rest-contract-test/parser/raml/v08"
	"github.com/jancajthaml/rest-contract-test/parser/raml/v10"
)

func fillResponses(method *v08.Method) []model.Response {
	// FIXME TBD
	return nil
}

func fillRequest(method *v08.Method) model.Request {
	// FIXME TBD
	return model.Request{}
}

func appendEndpoint(contract *model.Contract, path, method string) {
	res := model.Endpoint{
		Path:   path,
		Method: method,
		//Responses: fillResponses(endpoint),
		//Headers:   "",
		//Request:   fillRequest(endpoint),
	}

	contract.Endpoints = append(contract.Endpoints, res)
}

func extractMethods(contract *model.Contract, path string, resource *v08.Resource) {
	var foundSome = false

	if resource.Get != nil {
		appendEndpoint(contract, path, "GET")
		foundSome = true
	}

	if resource.Head != nil {
		appendEndpoint(contract, path, "HEAD")
		foundSome = true
	}

	if resource.Post != nil {
		appendEndpoint(contract, path, "POST")
		foundSome = true
	}

	if resource.Put != nil {
		appendEndpoint(contract, path, "PUT")
		foundSome = true
	}

	if resource.Patch != nil {
		appendEndpoint(contract, path, "PATCH")
		foundSome = true
	}

	if resource.Delete != nil {
		appendEndpoint(contract, path, "DELETE")
		foundSome = true
	}

	if foundSome {
		return
	}

	appendEndpoint(contract, path, "GET")
}

func walk(contract *model.Contract, path string, resource *v08.Resource) {
	extractMethods(contract, path, resource)

	for k, v := range resource.Nested {
		walk(contract, path+k, v)
	}
}

func FromFile(contract *model.Contract, file string) error {

	contract.Source = file

	switch io.GetDocumentType(file) {

	// INFO not implemented
	case "RAML 0.4":
		contract.Type = "RAML 0.4"

		_, err := v04.RAMLv04(file)
		if err != nil {
			return err
		}

		return nil

	case "RAML 0.8":
		contract.Type = "RAML 0.8"

		rootResource, err := v08.RAMLv08(file)
		if err != nil {
			return err
		}

		contract.Name = rootResource.Title

		for path, v := range rootResource.Resources {
			walk(contract, path, &v)
		}

		return nil

	// INFO not implemented
	case "RAML 1.0":
		contract.Type = "RAML 1.0"

		_, err := v10.RAMLv10(file)
		if err != nil {
			return err
		}

		return nil

	default:
		contract.Type = "Invalid"

		return fmt.Errorf("unsupported document")
	}

}
