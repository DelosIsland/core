// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package node
import (
	"fmt"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/DelosIsland/core/dngine"
	"github.com/DelosIsland/core/dngine/types"
	"github.com/DelosIsland/core/module/lib/go-config"
	client "github.com/DelosIsland/core/module/lib/go-rpc/client"
)

type MyNode struct {
	running int64

	logger *zap.Logger

	Dngine      *dngine.Dngine
	DngineTune  *dngine.DngineTunes
	Application types.Application
	GenesisDoc  *types.GenesisDoc
}

func NewMyNode(logger *zap.Logger, appName string, conf config.Config) *MyNode {
	if _, ok := Apps[appName]; !ok {
		return nil
	}
	app := Apps[appName](conf)
	tune := &dngine.DngineTunes{Conf: conf.(*config.MapConfig)}
	engine := dngine.NewDngine(tune)
	engine.ConnectApp(app)
	shard := &MyNode{
		logger: logger,

		Application: app,
		Dngine:      engine,
		DngineTune:  tune,
		GenesisDoc:  engine.Genesis(),
	}

	engine.SetSpecialVoteRPC(shard.GetSpecialVote)

	return shard
}

func (s *MyNode) Start() error {
	if atomic.CompareAndSwapInt64(&s.running, 0, 1) {
		s.Application.Start()
		return s.Dngine.Start()
	}
	return fmt.Errorf("already started")
}

func (s *MyNode) Stop() bool {
	if atomic.CompareAndSwapInt64(&s.running, 1, 0) {
		s.Application.Stop()
		return s.Dngine.Stop()
	}
	return false
}

func (s *MyNode) IsRunning() bool {
	return atomic.LoadInt64(&s.running) == 1
}

func (s *MyNode) GetSpecialVote(data []byte, validator *types.Validator) ([]byte, error) {
	clientJSON := client.NewClientJSONRPC(s.logger, validator.RPCAddress) // all shard nodes share the same rpc address of the Node
	tmResult := new(types.RPCResult)
	_, err := clientJSON.Call("vote_special_op", []interface{}{s.GenesisDoc.ChainID, data}, tmResult)
	if err != nil {
		return nil, err
	}
	res := (*tmResult).(*types.ResultRequestSpecialOP)
	if res.Code == types.CodeType_OK {
		return res.Data, nil
	}
	return nil, fmt.Errorf(res.Log)
}
