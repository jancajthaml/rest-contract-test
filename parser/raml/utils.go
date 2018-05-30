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
	"bufio"
	"bytes"
	"strconv"

	"encoding/json"
	"encoding/xml"

	yaml "github.com/advance512/yaml" // INFO .regex support
	//yaml "gopkg.in/yaml.v2"

	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	gio "github.com/jancajthaml/rest-contract-test/io"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func untypedConvert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = untypedConvert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = untypedConvert(v)
		}
	}
	return i
}

func ReadFileContents(filePath string) ([]byte, error) {

	// FIXME better and faster

	if len(filePath) == 0 {
		return nil, fmt.Errorf("File cannot be nil: %s", filePath)
	}

	// FIXME faster file read
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil,
			fmt.Errorf("Could not read file %s (Error: %s)", filePath, err.Error())
	}

	if gio.IsJSON(filePath, data) {

		var body interface{}
		if err := json.Unmarshal(data, &body); err != nil {
			return nil, err
		}

		body = untypedConvert(body)
		if b, err := yaml.Marshal(body); err != nil {
			return nil, err
		} else {
			b = append([]byte("\n"), b...)
			return b, nil
		}
	}

	if gio.IsXML(filePath, data) {
		var body interface{}
		if err := xml.Unmarshal(data, &body); err != nil {
			return nil, err
		}

		body = untypedConvert(body)
		if b, err := yaml.Marshal(body); err != nil {
			return nil, err
		} else {
			b = append([]byte("\n"), b...)
			return b, nil
		}
	}

	return data, nil
}

func CopyMap(ref map[string]string) map[string]string {
	clone := make(map[string]string)
	for k, v := range ref {
		clone[k] = v
	}
	return clone
}

func populateSecurityQueryParams(dataset map[string]SecurityScheme) map[string]map[string]string {
	result := make(map[string]map[string]string)

	for k, v := range dataset {

		if v.DescribedBy.QueryParameters != nil {

			placeholder := make(map[string]string)

			for name, parameter := range v.DescribedBy.QueryParameters.Data {
				if parameter.Example != nil {
					switch typed := parameter.Example.(type) {
					case string:
						placeholder[name] = strings.Replace(typed, "\n", "", -1)
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
				} else if parameter.Enum != nil {
					placeholder[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
				} else if parameter.Type != nil {
					// FIXME now need to generate value based by validations and type
					switch typed := parameter.Type.(type) {
					case string:
						placeholder[name] = typed
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
				}
			}
			if len(placeholder) != 0 {
				result[k] = placeholder
			}
		}
	}

	return result
}

func populateTraitQueryParams(dataset map[string]*Trait) map[string]map[string]string {

	result := make(map[string]map[string]string)

	for k, v := range dataset {

		if v.QueryParameters != nil {

			placeholder := make(map[string]string)

			for name, parameter := range v.QueryParameters.Data {
				if parameter.Example != nil {
					switch typed := parameter.Example.(type) {
					case string:
						placeholder[name] = strings.Replace(typed, "\n", "", -1)
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
				} else if parameter.Enum != nil {
					placeholder[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
				} else if parameter.Type != nil {
					switch typed := parameter.Type.(type) {
					case string:
						placeholder[name] = typed
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
					// FIXME now need to generate value based by validations and type
				}
			}
			if len(placeholder) != 0 {
				result[k] = placeholder
			}
		}
	}

	return result
}

func populateSecurityHeaders(dataset map[string]SecurityScheme) map[string]map[string]string {

	result := make(map[string]map[string]string)

	//fmt.Printf("headers security", dataset)

	for k, v := range dataset {

		//fmt.Println("headers security", k, v.DescribedBy.Headers)

		if v.DescribedBy.Headers != nil {
			//fmt.Println(k, v.DescribedBy.Headers)

			placeholder := make(map[string]string)

			for name, parameter := range v.DescribedBy.Headers.Data {
				if parameter.Example != nil {
					switch typed := parameter.Example.(type) {
					case string:
						placeholder[name] = strings.Replace(typed, "\n", "", -1)
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
				} else if parameter.Enum != nil {
					placeholder[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
				} else if parameter.Type != nil {
					switch typed := parameter.Type.(type) {
					case string:
						placeholder[name] = typed
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
					// FIXME now need to generate value based by validations and type
				}
			}
			if len(placeholder) != 0 {
				result[k] = placeholder
			}

		}
	}

	return result
}

func populateTraitHeaders(dataset map[string]*Trait) map[string]map[string]string {

	result := make(map[string]map[string]string)

	for k, v := range dataset {

		if v.Headers != nil {
			//			fmt.Println(k, v.Headers)

			placeholder := make(map[string]string)

			for name, parameter := range v.Headers.Data {
				if parameter.Example != nil {
					switch typed := parameter.Example.(type) {
					case string:
						placeholder[name] = strings.Replace(typed, "\n", "", -1)
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
				} else if parameter.Enum != nil {
					placeholder[name] = parameter.Enum[rand.Intn(len(parameter.Enum)-1)]
				} else if parameter.Type != nil {
					switch typed := parameter.Type.(type) {
					case string:
						placeholder[name] = typed
					case int:
						placeholder[name] = strconv.Itoa(typed)
					}
					// FIXME now need to generate value based by validations and type
				}
			}
			if len(placeholder) != 0 {
				result[k] = placeholder
			}

		}
	}

	return result
}

func PreProcess(originalContents io.Reader, workingDirectory string) ([]byte, error) {

	var preprocessedContents bytes.Buffer

	scanner := bufio.NewScanner(originalContents)
	var line string

	for scanner.Scan() {
		line = scanner.Text()

		// FIXME better
		if idx := strings.Index(line, "!include"); idx != -1 {

			includeLength := len("!include ")

			includedFile := line[idx+includeLength:]

			preprocessedContents.Write([]byte(line[:idx]))

			includedContents, err := ReadFileContents(filepath.Join(workingDirectory, includedFile))

			if err != nil {
				return nil,
					fmt.Errorf("Error including file %s:\n    %s",
						includedFile, err.Error())
			}

			internalScanner := bufio.NewScanner(bytes.NewBuffer(includedContents))

			firstLine := true
			indentationString := ""

			for internalScanner.Scan() {
				internalLine := internalScanner.Text()

				preprocessedContents.WriteString(indentationString)
				if firstLine {
					indentationString = strings.Repeat(" ", idx)
					firstLine = false
				}

				preprocessedContents.WriteString(internalLine)
				preprocessedContents.WriteByte('\n')
			}

		} else {
			preprocessedContents.WriteString(line)
			preprocessedContents.WriteByte('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %s", err.Error())
	}

	return preprocessedContents.Bytes(), nil
}
