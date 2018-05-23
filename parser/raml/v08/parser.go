package v08

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "github.com/advance512/yaml"
)

func ParseFile(filePath string) (*APIDefinition, error) {

	workingDirectory, fileName := filepath.Split(filePath)

	mainFileBytes, err := readFileContents(workingDirectory, fileName)

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

		if ramlVersion != "#%RAML 0.8" {
			return nil, errors.New("Input file is not a RAML 0.8 file. Make " +
				"sure the file starts with #%RAML 0.8")
		}
	}

	preprocessedContentsBytes, err := preProcess(mainFileBuffer, workingDirectory)

	if err != nil {
		return nil, fmt.Errorf("Error preprocessing RAML file (Error: %s)", err.Error())
	}

	apiDefinition := new(APIDefinition)
	apiDefinition.RAMLVersion = ramlVersion

	err = yaml.Unmarshal(preprocessedContentsBytes, apiDefinition)
	if err != nil {
		ramlError := new(RamlError)

		if yamlErrors, ok := err.(*yaml.TypeError); ok {
			populateRAMLError(ramlError, yamlErrors)
		} else {
			ramlError.Errors = append(ramlError.Errors, err.Error())
		}

		return nil, ramlError
	}

	return apiDefinition, nil
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

func readFileContents(workingDirectory string, fileName string) ([]byte, error) {

	// FIXME better and faster

	filePath := filepath.Join(workingDirectory, fileName)

	if fileName == "" {
		return nil, fmt.Errorf("File name cannot be nil: %s", filePath)
	}

	fileContentsArray, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil,
			fmt.Errorf("Could not read file %s (Error: %s)",
				filePath, err.Error())
	}

	if fileContentsArray[0] == byte('[') || fileContentsArray[0] == '{' {

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

func preProcess(originalContents io.Reader, workingDirectory string) ([]byte, error) {

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

			includedContents, err := readFileContents(workingDirectory, includedFile)

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

func RAMLv08(file string) (*APIDefinition, error) {
	apiDefinition, err := ParseFile(file)
	if err != nil {
		return nil, err
	}

	return apiDefinition, nil
}
