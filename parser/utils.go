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
	"os"
	"strings"
	"io"
)

func GetDocumentType(file string) (string) {
	// FIXME check suffix

	// FIXME RAML below
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

	switch strings.Split(string(buf), "\n")[0] { // FIXME

	case "#%RAML 0.8":
		return "RAML 0.8"

	case "#%RAML 1.0":
		return "RAML 1.0"

	default:
		return ""
	}
}
