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

package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/sirupsen/logrus"

	"strings"

	ramlv10 "github.com/tsaikd/go-raml-parser/parser" // needs GCC (cgo) :/ :( !!!
	ramlv08 "gopkg.in/raml.v0"
)

func ReadFirstLine(absPath string) (bool, string) {
	f, err := os.OpenFile(absPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false, ""
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false, ""
	}

	size := fi.Size()
	if size > 20 {
		size = 20
	}
	buf := make([]byte, size)
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return false, ""
	}

	return true, strings.Split(string(buf), "\n")[0] // FIXME
}

func CmdTest(c *cli.Context) error {
	file := c.Args().First()
	if len(file) == 0 {
		return fmt.Errorf("no file provided")
	}

	ok, firstLine := ReadFirstLine(file)
	if !ok {
		return fmt.Errorf("cannot read file")
	}

	switch firstLine {

	// INFO does not work with includes but good MVP for now
	case "#%RAML 0.8":
		apiDefinition, err := RAMLv08(file)
		if err != nil {
			return err
		}

		fmt.Printf("+------------------------------------------------------------------------\n")
		fmt.Printf("| RAML %s\n", file)
		fmt.Printf("+------------------------------------------------------------------------\n")
		fmt.Printf("| title: %s\n", apiDefinition.Title)
		fmt.Printf("+------------------------------------------------------------------------\n")

		// Iterate and print all endpoints
		for k, v := range apiDefinition.Resources {
			if v.Get != nil {
				fmt.Printf("| GET     | %s\n", k)
			}
			if v.Head != nil {
				fmt.Printf("| HEAD    | %s\n", k)
			}
			if v.Post != nil {
				fmt.Printf("| POST    | %s\n", k)
			}
			if v.Put != nil {
				fmt.Printf("| PUT     | %s\n", k)
			}
			if v.Patch != nil {
				fmt.Printf("| PATCH   | %s\n", k)
			}
			if v.Delete != nil {
				fmt.Printf("| DELETE  | %s\n", k)
			}
		}

		fmt.Printf("+------------------------------------------------------------------------\n")

	// INFO Does not work
	case "#%RAML 1.0":
		apiDefinition, err := RAMLv10(file)
		if err != nil {
			return err
		}

		// Iterate and print all endpoints
		for _, v := range apiDefinition.Resources {
			logrus.Info(v)
		}

	default:
		return fmt.Errorf("unsupported version of RAML ")
	}

	return nil
}

func RAMLv08(file string) (*ramlv08.APIDefinition, error) {
	apiDefinition, err := ramlv08.ParseFile(file)
	if err != nil {
		return nil, err
	}

	return apiDefinition, nil
}

func RAMLv10(file string) (*ramlv10.RootDocument, error) {
	apiDefinition, err := ramlv10.NewParser().ParseFile(file)
	if err != nil {
		return nil, err
	}

	logrus.Info(apiDefinition)

	return &apiDefinition, nil
}
