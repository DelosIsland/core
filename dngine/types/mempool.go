// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package types

type IMempool interface {
	Lock()
	Unlock()
	Update(height int64, txs []Tx)
}
