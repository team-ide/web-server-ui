package main

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"web-server-ui/commons"
	"web-server-ui/servers"
)

func main() {

	for _, v := range os.Args {
		if v == "-version" || v == "-v" {
			commons.OutVersion()
			return
		}
	}

	wait := &sync.WaitGroup{}
	wait.Add(1)

	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		panic(err)
	}
	rootDir = filepath.ToSlash(rootDir)
	if !strings.HasSuffix(rootDir, "/") {
		rootDir += "/"
	}

	config := servers.Config{
		Host:    "0.0.0.0",
		Port:    11030,
		Context: "",
		DistDir: rootDir + "dist/",
	}

	server, err := servers.New(config)
	if err != nil {
		util.Logger.Error("server new error", zap.Any("config", config), zap.Error(err))
		return
	}

	err = server.Startup(func() {
		wait.Done()
	})
	if err != nil {
		util.Logger.Error("server startup error", zap.Any("config", config), zap.Error(err))
		return
	}
	wait.Wait()
}
