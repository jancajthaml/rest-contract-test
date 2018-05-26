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

package common

import (
	"bufio"
	"bytes"

	"encoding/json"

	yaml "github.com/advance512/yaml" // INFO .regex support
	//yaml "gopkg.in/yaml.v2"

	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	gio "github.com/jancajthaml/rest-contract-test/io"
)

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
	fileContentsArray, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil,
			fmt.Errorf("Could not read file %s (Error: %s)",
				filePath, err.Error())
	}

	// FIXME now determine if source file is json or yaml, currently working expecting yaml

	// INFO test if included file is a json maybe separate to its own function

	//fmt.Printf("%s is json? %t\n", filePath, gio.IsJSON(fileContentsArray))

	if gio.IsJSON(filePath, fileContentsArray) {

		var body interface{}
		if err := json.Unmarshal(fileContentsArray, &body); err != nil {
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

	return fileContentsArray, nil
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
