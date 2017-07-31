// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package state

import (
	"github.com/DelosIsland/core/module/lib/go-wire"
	"github.com/DelosIsland/core/dngine/plugin"
)

type (

	// Plugin defines the behavior of the core plugins
	IPlugin interface {

		// DeliverTx return false means the tx won't be pass on to proxy app
		DeliverTx(tx []byte, i int) (bool, error)

		// CheckTx return false means the tx won't be pass on to proxy app
		CheckTx(tx []byte) (bool, error)

		// BeginBlock just mock the abci Blockaware interface
		BeginBlock(*plugin.BeginBlockParams) (*plugin.BeginBlockReturns, error)

		// EndBlock just mock the abci Blockaware interface
		EndBlock(*plugin.EndBlockParams) (*plugin.EndBlockReturns, error)

		// Reset is called when u don't need to maintain the plugin status
		Reset()

		// InitPlugin custom the initialization of the plugin
		InitPlugin(*plugin.InitPluginParams)
	}
)

var (
	pluginTypeSpecialOP = byte(0x01)
)

func init() {
	_ = wire.RegisterInterface(
		struct{ IPlugin }{},
		wire.ConcreteType{&plugin.Specialop{}, pluginTypeSpecialOP},
	)
}
