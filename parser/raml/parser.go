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

	gio "github.com/jancajthaml/rest-contract-test/io"
)

func ParseFile(filePath string) (*APIDefinition, error) {

	mainFileBytes, err := gio.ReadLocalFile(filePath)
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

	queryParamsTraits := <-eventualQueryParamsTraits
	headersTraits := <-eventualHeadersTraits

	queryParamsSecurity := <-eventualQueryParamsSecurity
	headersSecurity := <-eventualHeadersSecurity

	var prefix = ""
	if rootResource.BaseUri != nil && len(rootResource.BaseUri.Data) > 0 {
		prefix = rootResource.BaseUri.Data
	} else {
		prefix = "http://localhost:8080"
	}

	if strings.HasPrefix(prefix, "http") {
		// pass
	} else if len(rootResource.Protocols) > 0 {
		prefix = strings.ToLower(rootResource.Protocols[0]) + ":/" + prefix
	} else {
		prefix = "http:/" + prefix
	}

	for path, v := range rootResource.Resources {
		walk(contract, prefix+path, &v,
			make(map[string]string), make(map[string]string), make(map[int]model.Payload),
			queryParamsSecurity, queryParamsTraits,
			headersSecurity, headersTraits)
	}

	contract.Type = rootResource.RAMLVersion

	return contract, nil
}

// FIXME optimise method signature
func processMethod(contract *model.Contract, path string, kind string, method *Method,
	queryStrings map[string]string, headers map[string]string, responses map[int]model.Payload,
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
				fmt.Println("processing header", name, parameter.Example)

				switch typed := parameter.Example.(type) {
				case string:
					headers[name] = strings.Replace(typed, "\n", "", -1)
				case int:
					headers[name] = strconv.Itoa(typed)
				case map[interface{}]interface{}:
				    for k := range typed {
				        headers[name] = "{"+k.(string)+"}"
				        break
				    }
				}

			} else if parameter.Enum != nil {
				headers[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
			} else if len(parameter.Type) != 0 {
				headers[name] = RandValue(parameter.Type)
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
				case map[interface{}]interface{}:
				    for k := range typed {
				        queryStrings[name] = "{"+k.(string)+"}"
				        break
				    }
				}
			} else if parameter.Enum != nil {
				queryStrings[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
			} else if len(parameter.Type) != 0 {
				queryStrings[name] = RandValue(parameter.Type)

			}
		}
	}

	// FIXME copy from responses
	rs := make(map[int]model.Payload, 0)
	if len(method.Responses) != 0 {
		//fmt.Println("checking (0)", kind, path)

		for code, response := range method.Responses {
			if response.Referenced != nil {
				fmt.Println("response is referenced")
				continue
			}

			bodies := processBodies(nil, response.Bodies)
			for _, payload := range bodies {
				rs[code] = model.Payload{
					Content: &payload,
				}
			}

			// FIXME missing headers in responses

		}
	}
	if method.Bodies == nil {
		contract.Endpoints = append(contract.Endpoints, &model.Endpoint{
			URI:          path,
			Method:       kind,
			QueryStrings: queryStrings,
			Request: model.Payload{
				Headers: headers,
			},
			Responses: rs,
		})
		return
	}

	bodies := processBodies(nil, method.Bodies)
	for _, payload := range bodies {
		contract.Endpoints = append(contract.Endpoints, &model.Endpoint{
			URI:          path,
			Method:       kind,
			QueryStrings: queryStrings,
			Request: model.Payload{
				Headers: headers,
				Content: &payload,
			},
			Responses: rs,
		})
	}

	return
}

func processBodies(types *ResourceTypes, bodies *Bodies) []model.Content {

	result := make([]model.Content, 0)

	if bodies.Referenced != nil {
		fmt.Println("body by reference", *bodies.Referenced)
		return result
	}

	for mime, body := range bodies.ForMIMEType {

		if body.Example != nil {
			result = append(result, model.Content{
				Example: gio.UntypedConvert(body.Example),
				Type:    mime,
			})
		}
	}

	return result
}

// FIXME optimise method signature
func walk(contract *model.Contract, path string, resource *Resource,
	queryStrings map[string]string, headers map[string]string, responses map[int]model.Payload,
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
			CopyMap(queryStrings), CopyMap(headers), responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Head != nil {
		processMethod(contract, path, "HEAD", resource.Head,
			CopyMap(queryStrings), CopyMap(headers), responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Post != nil {
		processMethod(contract, path, "POST", resource.Post,
			CopyMap(queryStrings), CopyMap(headers), responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Put != nil {
		processMethod(contract, path, "PUT", resource.Put,
			CopyMap(queryStrings), CopyMap(headers), responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Patch != nil {
		processMethod(contract, path, "PATCH", resource.Patch,
			CopyMap(queryStrings), CopyMap(headers), responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if resource.Delete != nil {
		processMethod(contract, path, "DELETE", resource.Delete,
			CopyMap(queryStrings), CopyMap(headers), responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
		found = true
	}

	if !found {
		qs = CopyMap(queryStrings)
		hds = CopyMap(headers)
		contract.Endpoints = append(contract.Endpoints, &model.Endpoint{
			URI:          path,
			Method:       "GET",
			QueryStrings: qs,
			Request: model.Payload{
				Headers: hds,
			},
			Responses: responses,
		})
	}

	for k, v := range resource.Nested {
		walk(contract, path+k, v, queryStrings, headers, responses,
			securityQueryStrings, traitsQueryStrings,
			securityHeaders, traitsHeaders)
	}
}
