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
		return nil, err
	}

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

	eventualQueryParamsSecurity := make(chan map[string]map[string]string)
	eventualQueryParamsTraits := make(chan map[string]map[string]string)
	eventualHeadersTraits := make(chan map[string]map[string]string)
	eventualHeadersSecurity := make(chan map[string]map[string]string)

	// from securitySchemes
	go func() {
		eventualQueryParamsSecurity <- populateSecurityQueryParams(rootResource.SecuritySchemes)
	}()
	go func() {
		eventualHeadersSecurity <- populateSecurityHeaders(rootResource.SecuritySchemes)
	}()
	// from traits
	go func() {
		if rootResource.Traits != nil {
			eventualQueryParamsTraits <- populateTraitQueryParams(rootResource.Traits.Data)
		} else {
			eventualQueryParamsTraits <- make(map[string]map[string]string)
		}
	}()
	go func() {
		if rootResource.Traits != nil {
			eventualHeadersTraits <- populateTraitHeaders(rootResource.Traits.Data)
		} else {
			eventualHeadersTraits <- make(map[string]map[string]string)
		}
	}()

	// wait foall
	queryParamsSecurity := <-eventualQueryParamsSecurity
	queryParamsTraits := <-eventualQueryParamsTraits
	headersTraits := <-eventualHeadersTraits
	headersSecurity := <-eventualHeadersSecurity

	fmt.Println("headers:")
	fmt.Println("> security", headersSecurity)
	fmt.Println("> traits", headersTraits)

	for path, v := range rootResource.Resources {
		walk(contract, path, &v,
			make(map[string]string), make(map[string]string),
			queryParamsSecurity, queryParamsTraits,
			headersSecurity, headersTraits)
	}

	contract.Type = rootResource.RAMLVersion

	return contract, nil
}

func walk(contract *model.Contract, path string, resource *Resource,
	queryStrings map[string]string, headers map[string]string,
	securityQueryStrings map[string]map[string]string, traitsQueryStrings map[string]map[string]string,
	securityHeaders map[string]map[string]string, traitsHeaders map[string]map[string]string) {

	var found = false
	var qs map[string]string
	var hds map[string]string

	if resource.Is != nil {
		for _, ref := range resource.Is.Data {
			if val, ok := traitsQueryStrings[ref]; ok {
				for k, v := range val {
					queryStrings[k] = v
				}
			}
			if val, ok := traitsHeaders[ref]; ok {
				for k, v := range val {
					headers[k] = v
				}
			}
		}
	}

	if resource.Get != nil {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		if resource.Get.Is != nil {
			for _, ref := range resource.Get.Is.Data {
				if val, ok := traitsQueryStrings[ref]; ok {
					for k, v := range val {
						qs[k] = v
					}
				}
				if val, ok := traitsHeaders[ref]; ok {
					for k, v := range val {
						hds[k] = v
					}
				}
			}
		}
		// FIXME there can be queryParams inlined
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "GET",
			QueryStrings: qs,
			Headers:      hds,
		})
		found = true
	}

	if resource.Head != nil {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		if resource.Head.Is != nil {
			for _, ref := range resource.Head.Is.Data {
				if val, ok := traitsQueryStrings[ref]; ok {
					for k, v := range val {
						qs[k] = v
					}
				}
				if val, ok := traitsHeaders[ref]; ok {
					for k, v := range val {
						hds[k] = v
					}
				}
			}
		}
		// FIXME there can be queryParams inlined
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "HEAD",
			QueryStrings: qs,
			Headers:      hds,
		})
		found = true
	}

	if resource.Post != nil {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		if resource.Post.Is != nil {
			for _, ref := range resource.Post.Is.Data {
				if val, ok := traitsQueryStrings[ref]; ok {
					for k, v := range val {
						qs[k] = v
					}
				}
				if val, ok := traitsHeaders[ref]; ok {
					for k, v := range val {
						hds[k] = v
					}
				}
			}
		}
		// FIXME there can be queryParams inlined
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "POST",
			QueryStrings: qs,
			Headers:      hds,
		})
		found = true
	}

	if resource.Put != nil {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		if resource.Put.Is != nil {
			for _, ref := range resource.Put.Is.Data {
				if val, ok := traitsQueryStrings[ref]; ok {
					for k, v := range val {
						qs[k] = v
					}
				}
				if val, ok := traitsHeaders[ref]; ok {
					for k, v := range val {
						hds[k] = v
					}
				}
			}
		}
		// FIXME there can be queryParams inlined
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "PUT",
			QueryStrings: qs,
			Headers:      hds,
		})
		found = true
	}

	if resource.Patch != nil {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		if resource.Patch.Is != nil {
			for _, ref := range resource.Patch.Is.Data {
				if val, ok := traitsQueryStrings[ref]; ok {
					for k, v := range val {
						qs[k] = v
					}
				}
				if val, ok := traitsHeaders[ref]; ok {
					for k, v := range val {
						hds[k] = v
					}
				}
			}
		}
		// FIXME there can be queryParams inlined
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "PATCH",
			QueryStrings: qs,
			Headers:      hds,
		})
		found = true
	}

	if resource.Delete != nil {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		if resource.Delete.Is != nil {
			for _, ref := range resource.Delete.Is.Data {
				if val, ok := traitsQueryStrings[ref]; ok {
					for k, v := range val {
						qs[k] = v
					}
				}
				if val, ok := traitsHeaders[ref]; ok {
					for k, v := range val {
						hds[k] = v
					}
				}
			}
		}
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "DELETE",
			QueryStrings: qs,
			Headers:      hds,
		})
		found = true
	}

	if !found {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		contract.Endpoints = append(contract.Endpoints, model.Endpoint{
			Path:         path,
			Method:       "GET",
			QueryStrings: qs,
			Headers:      hds,
		})
	}

	for k, v := range resource.Nested {
		walk(contract, path+k, v, queryStrings, headers,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
	}
}
