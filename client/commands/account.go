// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package client


import (
	"gopkg.in/urfave/cli.v1"
	cl "github.com/DelosIsland/core/module/lib/go-rpc/client"
	"encoding/json"
)

var (
	AccountCommands = cli.Command{
		Name:     "account",
		Usage:    "operations for account",
		Category: "Account",
		Subcommands: []cli.Command{
			{
				Name:     "add",
				Action:   addAccount,
				Usage:    "create a new account",
				Category: "Account",
			},
			{
				Name:     "getAccount",
				Action:   getAccountAmount,
				Usage:    "create a new account",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "addr",
						Usage: "address",
					},
				},
			},
		},
	}
)

func addAccount() {

}
//./client --backend "tcp://localhost:46657" account getAccount --addr=source1
func getAccountAmount(ctx *cli.Context) error {
	addr := ctx.String("addr")
	clientJSON := cl.NewClientJSONRPC(logger, "tcp://localhost:46657")
	tmResult := &TmResult{}
	res, err := clientJSON.Call("get_account_amount", []interface{}{addr}, tmResult)
	println(res)
	rb,_ := json.Marshal(res)
	println(string(rb))
	println(tmResult.Amount)
	return err
}


type TmResult struct {
	Amount int `json:"amount"`
}