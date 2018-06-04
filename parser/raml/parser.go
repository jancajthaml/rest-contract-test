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
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

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

	for path, v := range rootResource.Resources {
		walk(contract, path, &v,
			make(map[string]string), make(map[string]string),
			queryParamsSecurity, queryParamsTraits,
			headersSecurity, headersTraits)
	}

	contract.Type = rootResource.RAMLVersion

	return contract, nil
}

func processMethod(contract *model.Contract, path string, kind string, method *Method,
	queryStrings map[string]string, headers map[string]string,
	securityQueryStrings map[string]map[string]string, traitsQueryStrings map[string]map[string]string,
	securityHeaders map[string]map[string]string, traitsHeaders map[string]map[string]string) {

	if method.Is != nil {
		for _, ref := range method.Is.Data {
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

	if method.SecuredBy != nil {
		for _, ref := range method.SecuredBy.Data {
			if val, ok := securityQueryStrings[ref]; ok {
				for k, v := range val {
					queryStrings[k] = v
				}
			}
			if val, ok := securityHeaders[ref]; ok {
				for k, v := range val {
					headers[k] = v
				}
			}
		}
	}

	if method.Headers != nil {
		for name, parameter := range method.Headers.Data {
			if parameter.Example != nil {
				switch typed := parameter.Example.(type) {
				case string:
					headers[name] = strings.Replace(typed, "\n", "", -1)
				case int:
					headers[name] = strconv.Itoa(typed)
				}
			} else if parameter.Enum != nil {
				headers[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
			} else if parameter.Type != nil {
				switch typed := parameter.Type.(type) {
				case string:
					headers[name] = typed
				case int:
					headers[name] = strconv.Itoa(typed)
				}
				// FIXME now need to generate value based by validations and type
			}
		}
	}

	if method.QueryParameters != nil {
		for name, parameter := range method.QueryParameters.Data {
			if parameter.Example != nil {
				switch typed := parameter.Example.(type) {
				case string:
					queryStrings[name] = strings.Replace(typed, "\n", "", -1)
				case int:
					queryStrings[name] = strconv.Itoa(typed)
				}
			} else if parameter.Enum != nil {
				queryStrings[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
			} else if parameter.Type != nil {
				switch typed := parameter.Type.(type) {
				case string:
					queryStrings[name] = typed
				case int:
					queryStrings[name] = strconv.Itoa(typed)
				}
				// FIXME now need to generate value based by validations and type
			}
		}
	}

	contract.Endpoints = append(contract.Endpoints, model.Endpoint{
		Path:         path,
		Method:       kind,
		QueryStrings: queryStrings,
		Headers:      headers,
	})
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

	if resource.SecuredBy != nil {
		for _, ref := range resource.SecuredBy.Data {
			if val, ok := securityQueryStrings[ref]; ok {
				for k, v := range val {
					queryStrings[k] = v
				}
			}
			if val, ok := securityHeaders[ref]; ok {
				for k, v := range val {
					headers[k] = v
				}
			}
		}
	}

	if resource.Get != nil {
		processMethod(contract, path, "GET", resource.Get,
			CopyMap(queryStrings), CopyMap(headers),
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Head != nil {
		processMethod(contract, path, "HEAD", resource.Head,
			CopyMap(queryStrings), CopyMap(headers),
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Post != nil {
		processMethod(contract, path, "POST", resource.Post,
			CopyMap(queryStrings), CopyMap(headers),
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Put != nil {
		processMethod(contract, path, "PUT", resource.Put,
			CopyMap(queryStrings), CopyMap(headers),
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Patch != nil {
		processMethod(contract, path, "PATCH", resource.Patch,
			CopyMap(queryStrings), CopyMap(headers),
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Delete != nil {
		processMethod(contract, path, "DELETE", resource.Delete,
			CopyMap(queryStrings), CopyMap(headers),
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
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
