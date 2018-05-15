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

	"github.com/codegangsta/cli"
	"github.com/sirupsen/logrus"
)

func CmdTest(c *cli.Context) error {
	file := c.Args().First()
	if len(file) == 0 {
		return fmt.Errorf("no file provided")
	}

	logrus.Info("placeholder 1", file)
	logrus.Debug("placeholder 2")
	return nil
}
