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
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	//"net/http"

	gio "github.com/jancajthaml/rest-contract-test/io"
	"github.com/jancajthaml/rest-contract-test/model"
	"github.com/jancajthaml/rest-contract-test/parser/raml"
	"github.com/jancajthaml/rest-contract-test/parser/swagger"

	"io"
)

func FromResource(resource string) (*model.Contract, error) {
	if _, err := url.ParseRequestURI(resource); err == nil {
		return fromUri(resource)
	}

	return fromFile(resource)
}

func fromUri(uri string) (*model.Contract, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	buffer := bytes.NewBuffer(nil)

	if _, err := io.Copy(buffer, response.Body); err != nil {
		return nil, err
	}

	//fmt.Println(string(buffer.Bytes()))

	return nil, fmt.Errorf("loading from uri not implemented")
}

func fromFile(file string) (*model.Contract, error) {

	switch gio.GetDocumentType(file) {

	case "RAML":
		contract, err := raml.NewRaml(file)
		if err != nil {
			return nil, err
		}

		return contract, nil

	case "SWAGGER":
		contract, err := swagger.NewSwagger(file)
		if err != nil {
			return nil, err
		}

		return contract, nil

	default:
		return nil, fmt.Errorf("unsupported document")

	}
}
