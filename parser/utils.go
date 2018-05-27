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
	"github.com/jancajthaml/rest-contract-test/parser/raml"
)

func FromFile(file string) (*model.Contract, error) {

	switch io.GetDocumentType(file) {

	case "RAML 0.8":
		contract, err := raml.NewRAML(file)
		if err != nil {
			return nil, err
		}

		return contract, nil

	case "RAML 1.0":
		contract, err := raml.NewRAML(file)
		if err != nil {
			return nil, err
		}

		return contract, nil

	default:
		return nil, fmt.Errorf("unsupported document")

	}
}
