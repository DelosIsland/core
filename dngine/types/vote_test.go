// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package types

import (
	"testing"
)

func TestVoteSignable(t *testing.T) {
	vote := &Vote{
		ValidatorAddress: []byte("addr"),
		ValidatorIndex:   56789,
		Height:           12345,
		Round:            23456,
		Type:             byte(2),
		BlockID: BlockID{
			Hash: []byte("hash"),
			PartsHeader: PartSetHeader{
				Total: 1000000,
				Hash:  []byte("parts_hash"),
			},
		},
	}
	signBytes := SignBytes("test_chain_id", vote)
	signStr := string(signBytes)

	expected := `{"chain_id":"test_chain_id","vote":{"block_id":{"hash":"68617368","parts":{"hash":"70617274735F68617368","total":1000000}},"height":12345,"round":23456,"type":2}}`
	if signStr != expected {
		// NOTE: when this fails, you probably want to fix up consensus/replay_test too
		t.Errorf("Got unexpected sign string for Vote. Expected:\n%v\nGot:\n%v", expected, signStr)
	}
}
