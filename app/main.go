// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package main

import (
	"flag"
	"path"

	"github.com/DelosIsland/core/app/node"
	"github.com/DelosIsland/core/dngine"
	acfg "github.com/DelosIsland/core/dngine/config"
	"go.uber.org/zap"
)

// params
var (
	InitFlag = flag.Bool("init", false, "set initial files")
	DataDir  = flag.String("datadir", "", "set data direction")
)

func main() {
	flag.Parse()
	configPath := *DataDir //"/home/vagrant/gohome/src/github.com/DelosIsland/core/dngine-test/src/config2/"
	conf := acfg.GetConfig(configPath)
	env := conf.GetString("environment")
	logger := Initialize(env, path.Join("", "node.output.log"), path.Join("", "node.err.log"))
	if *InitFlag == true {
		dngine.Initialize(&dngine.DngineTunes{Conf: conf})
	} else {
		node.RunNode(logger, conf)
	}
}

// Initialize init log
func Initialize(env, output, errOutput string) *zap.Logger {
	var zapConf zap.Config
	var err error

	if env == "production" {
		zapConf = zap.NewProductionConfig()
	} else {
		zapConf = zap.NewDevelopmentConfig()
	}

	zapConf.OutputPaths = []string{output}
	zapConf.ErrorOutputPaths = []string{errOutput}
	logger, err := zapConf.Build()
	if err != nil {
		panic(err.Error())
	}

	logger.Debug("Starting zap! Have your fun!")

	return logger
}
