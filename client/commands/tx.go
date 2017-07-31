// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package client

import (
	"gopkg.in/urfave/cli.v1"
	cl "github.com/DelosIsland/core/module/lib/go-rpc/client"
	"encoding/json"
	"github.com/DelosIsland/core/app/node"
	"github.com/DelosIsland/core/dngine/types"
)

var (
	TxCommands = cli.Command{
		Name:     "tx",
		Usage:    "operations for transaction",
		Category: "Transaction",
		Subcommands: []cli.Command{
			{
				Name:   "send",
				Usage:  "send a transaction",
				Action: sendTx,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "addr",
						Usage: "address",
					},
					cli.StringFlag{
						Name:  "payload",
						Usage: "payload",
					},
					cli.StringFlag{
						Name:  "privKey",
						Usage: "privateKey",
					},
					cli.StringFlag{
						Name:  "dest",
						Usage: "destAddress",
					},
					cli.StringFlag{
						Name:  "amount",
						Usage: "amount",
					},
					cli.StringFlag{
						Name:  "chainid",
						Usage: "chainid",
					},
				},
			},
		},
	}
)

// command params
// --backend "tcp://localhost:46657" tx send --addr=source1 --dest=dest1 --amount=1000 --chainid=dngine-test
func sendTx(ctx *cli.Context) error {
	destAddr := ctx.String("dest")
	sourceAddr := ctx.String("addr")
	amount := ctx.Int64("amount")
	chainID := ctx.String("chainid")

	myTx := node.MyTx{
		DestAddress:destAddr,
		SourceAddress: sourceAddr,
		Amount: int(amount),
	}

	txContent, _ := json.Marshal(myTx)
	txTotal := []byte{'a', 'b', 'c', 'd'}
	for _, v := range txContent {
		txTotal = append(txTotal, v)
	}

	clientJSON := cl.NewClientJSONRPC(logger, "tcp://localhost:46657")
	res := &types.ResultBroadcastTx{}
	_, err := clientJSON.Call("broadcast_tx_sync", []interface{}{chainID, txTotal}, res)

	return err
}