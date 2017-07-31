// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package types

import (
	"github.com/DelosIsland/core/module/lib/go-config"
)

type Application interface {
	GetDngineHooks() Hooks
	CompatibleWithDngine()
	CheckTx([]byte) error
	Query([]byte) Result
	Info() ResultInfo
	Start()
	Stop()
}

type AppMaker func(config.Config) Application
