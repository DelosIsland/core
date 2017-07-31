// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package xlog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	//"log"
	"runtime"
	"time"
)

func DumpStack() {
	if err := recover(); err != nil {
		var buf bytes.Buffer
		bs := make([]byte, 1<<12)
		num := runtime.Stack(bs, false)
		buf.WriteString(fmt.Sprintf("Panic: %s\n", err))
		buf.Write(bs[:num])
		dumpName := logDir + "/dump_" + time.Now().Format("20060102-150405")
		nerr := ioutil.WriteFile(dumpName, buf.Bytes(), 0644)
		if nerr != nil {
			fmt.Println("write dump file error", nerr)
			fmt.Println(buf.Bytes())
		}
	}
}
