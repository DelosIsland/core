// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package plugin

import (
	"github.com/DelosIsland/core/dngine/refuse_list"
	"github.com/DelosIsland/core/dngine/types"
	"github.com/DelosIsland/core/module/lib/go-crypto"
	"github.com/DelosIsland/core/module/lib/go-p2p"
)

const (
	PluginNoncePrefix = "pn-"
)

type (
	InitPluginParams struct {
		Switch     *p2p.Switch
		PrivKey    crypto.PrivKeyEd25519
		RefuseList *refuse_list.RefuseList
		Validators **types.ValidatorSet
	}

	BeginBlockParams struct {
		Block *types.Block
	}

	BeginBlockReturns struct {
	}

	EndBlockParams struct {
		Block             *types.Block
		ChangedValidators []*types.ValidatorAttr
		NextValidatorSet  *types.ValidatorSet
	}

	EndBlockReturns struct {
		NextValidatorSet *types.ValidatorSet
	}
)
