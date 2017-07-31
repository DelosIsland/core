// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package dngine

import (
	"path"

	"go.uber.org/zap"
)

func InitializeLog(env, logpath string) *zap.Logger {
	var zapConf zap.Config
	var err error

	if env == "production" {
		zapConf = zap.NewProductionConfig()
	} else {
		zapConf = zap.NewDevelopmentConfig()
	}

	zapConf.OutputPaths = []string{path.Join(logpath, "output.log")}
	zapConf.ErrorOutputPaths = []string{path.Join(logpath, "err.output.log")}
	logger, err := zapConf.Build()
	if err != nil {
		panic(err.Error())
	}

	return logger
}
