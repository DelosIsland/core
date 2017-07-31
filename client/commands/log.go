// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package client

import (
	"os"
	"path"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	zapConf := zap.NewDevelopmentConfig()
	zapConf.OutputPaths = []string{path.Join(pwd, "client.out.log")}
	zapConf.ErrorOutputPaths = []string{}
	logger, err = zapConf.Build()
	if err != nil {
		panic(err.Error())
	}
}
