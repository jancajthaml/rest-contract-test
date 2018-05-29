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

package raml

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jancajthaml/rest-contract-test/model"

	yaml "github.com/advance512/yaml" // INFO .regex support
	//yaml "gopkg.in/yaml.v2"
)

func ParseFile(filePath string) (*APIDefinition, error) {

	mainFileBytes, err := ReadFileContents(filePath)
	if err != nil {
		return nil, err
	}

	mainFileBuffer := bytes.NewBuffer(mainFileBytes)

	var ramlVersion string
	if firstLine, err := mainFileBuffer.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("Problem reading RAML file (Error: %s)", err.Error())
	} else {

		if len(firstLine) >= 10 {
			ramlVersion = firstLine[:10]
		}

		if ramlVersion != "#%RAML 1.0" && ramlVersion != "#%RAML 0.8" {
			return nil, errors.New("Resource is not RAML 0.8 or 1.0")
		}
	}

	workingDirectory, _ := filepath.Split(filePath)
	preprocessedContentsBytes, err := PreProcess(mainFileBuffer, workingDirectory)

	if err != nil {
		return nil, fmt.Errorf("Error preprocessing RAML file (Error: %s)", err.Error())
	}

	apiDefinition := new(APIDefinition)
	apiDefinition.RAMLVersion = ramlVersion[2:]

	err = yaml.Unmarshal(preprocessedContentsBytes, apiDefinition)
	if err != nil {
		fmt.Println(string(preprocessedContentsBytes))
		return nil, err
	}

	PostProcess(apiDefinition)

	return apiDefinition, nil
}

func NewRaml(file string) (*model.Contract, error) {
	contract := new(model.Contract)

	contract.Source = file

	rootResource, err := ParseFile(file)
	if err != nil {
		return contract, err
	}

	contract.Name = rootResource.Title

	for path, v := range rootResource.Resources {
		walk(contract, path, &v)
	}

	contract.Type = rootResource.RAMLVersion

	return contract, nil
}

func fillResponses(method *Method) []model.Response {
	// FIXME TBD
	return nil
}

func fillRequest(method *Method) model.Request {
	// FIXME TBD
	return model.Request{}
}

func extractMethods(contract *model.Contract, path string, resource *Resource) {
	var method = ""

	if resource.Get != nil {
		method = "GET"
		//fmt.Println(resource.Get.Headers)
	}

	if resource.Head != nil {
		method = "HEAD"
	}

	if resource.Post != nil {
		method = "POST"
	}

	if resource.Put != nil {
		method = "PUT"
	}

	if resource.Patch != nil {
		method = "PATCH"
	}

	if resource.Delete != nil {
		method = "DELETE"
	}

	if len(method) == 0 {
		method = "GET"
	}

	res := model.Endpoint{
		Path:   path,
		Method: method,
		//Responses: fillResponses(endpoint),
		//Headers:   "",
		//Request:   fillRequest(endpoint),
	}

	contract.Endpoints = append(contract.Endpoints, res)
}

func walk(contract *model.Contract, path string, resource *Resource) {
	extractMethods(contract, path, resource)

	for k, v := range resource.Nested {
		walk(contract, path+k, v)
	}
}
