// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
	"github.com/DelosIsland/core/client/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "tool"

	app.Commands = []cli.Command{
		client.AccountCommands,
		client.TxCommands,
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "backend",
			Value:       "tcp://localhost:46657",
		},
	}

	_ = app.Run(os.Args)
}
