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
	"os"
	"fmt"

	"github.com/codegangsta/cli"
)

func GlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose logging",
		},
	}
}

func All() []cli.Command {
	return []cli.Command{
		{
			Name:   "test",
			Usage:  "FILE, tests contract provided documentation file",
			Action: try(CmdTest),
		},
	}
}

func NotFound(c *cli.Context, command string) {
	cli.ShowAppHelp(c)
	os.Exit(2)
}

func try(fn func(c *cli.Context) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		// FIXME deffer error here

		if err := fn(c); err != nil {
			if c.GlobalBool("verbose") {
				panic(err)
			} else {
				fmt.Println(err)
				fmt.Println("command failed. use --verbose to see full stacktrace")
				return err
			}
		}
		return nil
	}
}
