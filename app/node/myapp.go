// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/DelosIsland/core/dngine/types"
	cmn "github.com/DelosIsland/core/module/lib/go-common"
	"github.com/DelosIsland/core/module/lib/go-config"
	"github.com/DelosIsland/core/module/lib/go-crypto"
	"github.com/DelosIsland/core/module/lib/go-db"
	"github.com/DelosIsland/core/module/lib/go-merkle"
	"github.com/DelosIsland/core/module/lib/go-wire"
)

//const (
//	ShardTag = "shrd"
//
//	ShardRegister byte = iota
//	ShardCreate
//	ShardJoin
//	ShardLeave
//	ShardDelete
//)

type LastBlockInfo struct {
	Height uint64
	Txs    [][]byte
	Hash   []byte
}

type ManagedState struct {
	accounts map[string]int
}

var managedState = ManagedState{accounts:make(map[string]int)}

type MyApp struct {
	config      config.Config
	privkey     crypto.PrivKeyEd25519
	chainDb     *db.GoLevelDB
	node        *Node
	engineHooks types.Hooks
	logger      *zap.Logger

	mtx sync.Mutex

	Txs             [][]byte
	RunningChainIDs map[string]string
	Shards          map[string]*MyNode
}

type Accounts struct {
	Address string
}

type MyTx struct {
	App             string                  `json:"app"`
	Act             byte                    `json:"act"`
	ChainID         string                  `json:"chainid"`
	Genesis         types.GenesisDoc        `json:"genesis"`
	Config          map[string]interface{}  `json:"config"`
	Time            time.Time               `json:"time"`
	Signature       crypto.SignatureEd25519 `json:"signature"`
	DestAddress	    string   		        `json:"destaddress"`
	SourceAddress	string		            `json:"sourceaddress"`
	Amount          int	                    `json:"amount"`
}

func (t *MyTx) SignByPrivKey(p crypto.PrivKeyEd25519) error {
	txBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	t.Signature = p.Sign(txBytes).(crypto.SignatureEd25519)

	return nil
}

func (t *MyTx) VerifySignature(pubKey crypto.PubKeyEd25519, signature crypto.SignatureEd25519) bool {
	//tBytes, err := json.Marshal(t)
	//if err != nil {
	//	return false
	//}
	//pubKeyBytes := [32]byte(pubKey)
	//sig := [64]byte(signature)
	//return ed25519.Verify(&pubKeyBytes, tBytes, &sig)
	return true
}

var (
	lastBlockKey = []byte("lastblock")

	ErrUnknownTx   = fmt.Errorf("unknown tx")
	ErrShardExists = fmt.Errorf("shard exists")
	ErrFailToStop  = fmt.Errorf("stop shard failed")
)

func NewMyApp(logger *zap.Logger, conf config.Config) *MyApp {
	datadir := conf.GetString("db_dir")
	app := MyApp{
		config: conf,
		logger: logger,

		Txs:             make([][]byte, 0),
		RunningChainIDs: make(map[string]string),
		Shards:          make(map[string]*MyNode),
	}

	var err error
	if app.chainDb, err = db.NewGoLevelDB("chaindata", datadir); err != nil {
		cmn.PanicCrisis(err)
	}

	app.engineHooks = types.Hooks{
		OnExecute: types.NewHook(app.OnExecute),
		OnCommit:  types.NewHook(app.OnCommit),
	}

	return &app
}

func (app *MyApp) Lock() {
	app.mtx.Lock()
}

func (app *MyApp) Unlock() {
	app.mtx.Unlock()
}

func (app *MyApp) setNode(n *Node) {
	app.node = n
}

// Stop stops all still running shards
func (app *MyApp) Stop() {
	for i := range app.Shards {
		app.Shards[i].Stop()
	}
}

// Start will restore all shards according to tx history
func (app *MyApp) Start() {
	lastBlock := app.LoadLastBlock()
	if lastBlock.Txs != nil && len(lastBlock.Txs) > 0 {
		app.Txs = lastBlock.Txs
		for _, tx := range lastBlock.Txs {
			app.ExecuteTx(nil, tx, 0)
		}
	}
}

func (app *MyApp) GetDngineHooks() types.Hooks {
	return app.engineHooks
}

func (app *MyApp) CompatibleWithDngine() {}

// ExecuteTx execute tx one by one in the loop, without lock, so should always be called between Lock() and Unlock() on the *stateDup
func (app *MyApp) ExecuteTx(blockHash []byte, bs []byte, txIndex int) (validtx []byte, err error) {
	//if !app.IsShardingTx(bs) {
	//	return nil, ErrUnknownTx
	//}

	txBytes := types.UnwrapTx(bs)
	tx := MyTx{}
	if err := json.Unmarshal(txBytes, &tx); err != nil {
		app.logger.Info("Unmarshal tx failed", zap.Binary("tx", txBytes), zap.String("error", err.Error()))
		return nil, err
	}

	sig := tx.Signature
	tx.Signature = [64]byte{}
	pubkey, ok := app.node.PrivValidator().PubKey.(crypto.PubKeyEd25519)
	if !ok {
		return nil, fmt.Errorf("node's pubkey must be crypto.PubKeyEd25519")
	}

	if !tx.VerifySignature(pubkey, sig) {
		app.logger.Debug("this tx is not for me", zap.Binary("tx", bs))
		return bs, nil
	}

	baseAmount := 1000000
	if _, ok := managedState.accounts[tx.SourceAddress]; !ok {
		managedState.accounts[tx.SourceAddress] = baseAmount
	}
	if _, ok := managedState.accounts[tx.DestAddress]; !ok {
		managedState.accounts[tx.DestAddress] = baseAmount
	}
	managedState.accounts[tx.SourceAddress] = managedState.accounts[tx.SourceAddress] - tx.Amount
	managedState.accounts[tx.DestAddress] = managedState.accounts[tx.DestAddress] + tx.Amount

	return bs, nil
}

// OnExecute would not care about Block.ExTxs
func (app *MyApp) OnExecute(height, round int, block *types.Block) (interface{}, error) {
	var (
		res types.ExecuteResult
		err error

		blockHash = block.Hash()
	)

	for i := range block.Txs {
		vtx, err := app.ExecuteTx(blockHash, block.Txs[i], i)
		if err == nil {
			res.ValidTxs = append(res.ValidTxs, vtx)
		} else {
			if err == ErrUnknownTx {
				// maybe we could do something with another app or so
			} else {
				res.InvalidTxs = append(res.InvalidTxs, types.ExecuteInvalidTx{Bytes: block.Txs[i], Error: err})
			}
		}
	}

	app.Txs = append(app.Txs, res.ValidTxs...)

	return res, err
}

// OnCommit run in a sync way, we don't need to lock stateDupMtx, but stateMtx is still needed
func (app *MyApp) OnCommit(height, round int, block *types.Block) (interface{}, error) {
	lastBlock := LastBlockInfo{Height: uint64(height), Txs: app.Txs, Hash: merkle.SimpleHashFromHashes(app.Txs)}
	app.SaveLastBlock(lastBlock)
	return types.CommitResult{AppHash: lastBlock.Hash}, nil
}

func (app *MyApp) LoadLastBlock() (lastBlock LastBlockInfo) {
	buf := app.chainDb.Get(lastBlockKey)
	if len(buf) != 0 {
		r, n, err := bytes.NewReader(buf), new(int), new(error)
		wire.ReadBinaryPtr(&lastBlock, r, 0, n, err)
		if *err != nil {

		}
	}
	return lastBlock
}

func (app *MyApp) SaveLastBlock(lastBlock LastBlockInfo) {
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(lastBlock, buf, n, err)
	if *err != nil {
		cmn.PanicCrisis(*err)
	}
	app.chainDb.SetSync(lastBlockKey, buf.Bytes())
}

func (app *MyApp) CheckTx(bs []byte) error {
	//if !app.IsShardingTx(bs) {
	//	return ErrUnknownTx
	//}

	txBytes := types.UnwrapTx(bs)
	tx := MyTx{}

	if err := json.Unmarshal(txBytes, &tx); err != nil {
		app.logger.Info("Unmarshal tx failed", zap.Binary("tx", txBytes), zap.String("error", err.Error()))
		return err
	}
	sig := tx.Signature
	tx.Signature = [64]byte{}
	pubkey, ok := app.node.PrivValidator().PubKey.(crypto.PubKeyEd25519)
	if !ok {
		return fmt.Errorf("my key is not crypto.PubKeyEd25519")
	}
	if !tx.VerifySignature(pubkey, sig) {
		return fmt.Errorf("wrong tx for wrong node: %X", pubkey)
	}

	return nil
}

func (app *MyApp) Info() (resInfo types.ResultInfo) {
	lb := app.LoadLastBlock()
	resInfo.LastBlockAppHash = lb.Hash
	resInfo.LastBlockHeight = lb.Height
	resInfo.Version = "alpha 0.1"
	resInfo.Data = "default app with sharding-controls"
	return
}

func (app *MyApp) Query(query []byte) types.Result {
	if len(query) == 0 {
		return types.NewResultOK([]byte{}, "Empty query")
	}
	var res types.Result
	action := query[0]
	// _ = query[1:]
	switch action {
	default:
		res = types.NewError(types.CodeType_BaseInvalidInput, "unimplemented query")
	}

	return res
}
