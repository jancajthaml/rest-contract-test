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

package io

import (
	"path/filepath"
	"strings"
)

func IsJSON(filename string, data []byte) bool {

	if strings.HasSuffix(filename, ".json") {
		return true
	}

	return data[0] == byte('{') || data[0] == byte('[')

	// FIXME not robust agains leading and trailing whitespace characters
	/*
		if data[0] == byte('[') && data[len(data)-1] == byte(']') {
			return true
		}

		if data[0] == byte('{') && data[len(data)-1] == byte('}') {
			return true
		}

		return false
	*/
}

func IsXML(filename string, data []byte) bool {

	if strings.HasSuffix(filename, ".xml") {
		return true
	}

	return data[0] == byte('<')

	// FIXME not robust agains leading and trailing whitespace characters
	/*
		if data[0] == byte('[') && data[len(data)-1] == byte(']') {
			return true
		}

		if data[0] == byte('{') && data[len(data)-1] == byte('}') {
			return true
		}

		return false
	*/
}

// FIXME convert to data []byte from local file
func GetDocumentType(file string) string {
	switch filepath.Ext(file) {

	case ".raml":
		// raml versions: { 0.8, 1.0, 2.0 }
		return "RAML" // + getRamlVersion(file)

	// https://github.com/yvasiyarov/swagger/tree/master/parser

	case ".json":

		// FIXME try to unmarshall swagger from json

		// swagger v2.0 header: swagger: 2.0
		// swagger v3.0 header: openapi: 3.0.0

		// https://github.com/OAI/OpenAPI-Specification
		// https://github.com/BigstickCarpet/swagger-parser
		// https://github.com/go-swagger/go-swagger
		// swagger versions: { 1.0, 1.1, 1.2, 2.0, 3.0, 3.1 }

		// FIXME swagger must be json check json
		return "SWAGGER"

	case ".yaml", ".yml":

		// FIXME try to unmarshall swagger from yml

		// swagger v2.0 header: swagger: 2.0
		// swagger v3.0 header: openapi: 3.0.0

		// https://github.com/OAI/OpenAPI-Specification
		// https://github.com/BigstickCarpet/swagger-parser
		// https://github.com/go-swagger/go-swagger
		// swagger versions: { 1.0, 1.1, 1.2, 2.0, 3.0, 3.1 }

		// FIXME swagger must be json check json
		return "SWAGGER"

	default:
		return ""

	}
}

/*
func getRamlVersion(resource string) string {

	// FIXME assuming that resource is local file
	file := resource

	// FIXME RAML from local file below
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return ""
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return ""
	}

	size := fi.Size()
	if size > 20 {
		size = 20
	}
	buf := make([]byte, size)
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return ""
	}

	switch strings.Split(string(buf), "\n")[0] { // FIXME do this better

	case "#%RAML 0.4":
		return "0.4"

	case "#%RAML 0.8":
		return "0.8"

	case "#%RAML 1.0":
		return "1.0"

	default:
		return ""
	}
}*/
