// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package consensus

import (
	. "github.com/DelosIsland/core/module/lib/go-common"
)

// kind of arbitrary
var Spec = "1"     // async
var Major = "0"    //
var Minor = "2"    // replay refactor
var Revision = "2" // validation -> commit

var Version = Fmt("v%s/%s.%s.%s", Spec, Major, Minor, Revision)
