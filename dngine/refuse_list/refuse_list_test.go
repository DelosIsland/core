// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package refuse_list

import (
	"os"
	"testing"

	"github.com/DelosIsland/core/dngine/types"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-db"
)

func TestRefuseList(t *testing.T) {
	refuseList := NewRefuseList(db.LevelDBBackendStr, "./")
	defer func() {
		os.RemoveAll("refuse_list.db")
		refuseList.db.Close()
	}()
	var keyStr = "6FEBD39916627AA0CD7CFDA4A94586F3BA958078621E6E466488A423272B9700"

	pubKey, err := types.StringTo32byte(keyStr)
	assert.Nil(t, err)
	refuseList.AddRefuseKey(pubKey)
	assert.Equal(t, true, refuseList.QueryRefuseKey(pubKey))
	assert.Equal(t, []string{keyStr}, refuseList.ListAllKey())
	refuseList.DeleteRefuseKey(pubKey)
	assert.Equal(t, 0, len(refuseList.ListAllKey()))
}
