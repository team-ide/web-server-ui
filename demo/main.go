package main

import (
	"fmt"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/web-server-ui/servers"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	version = "0.0.0"
)

func Version() string {
	return version
}

func OutVersion() {
	fmt.Println("Server Version : " + Version())
	fmt.Println("GO OS          : " + runtime.GOOS)
	fmt.Println("GO ARCH        : " + runtime.GOARCH)
	fmt.Println("GO Version     : " + runtime.Version())
}

func main() {

	for _, v := range os.Args {
		if v == "-version" || v == "-v" {
			OutVersion()
			return
		}
	}

	wait := &sync.WaitGroup{}
	wait.Add(1)

	localDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	localDir, err = filepath.Abs(localDir)
	if err != nil {
		panic(err)
	}
	localDir = filepath.ToSlash(localDir)
	if !strings.HasSuffix(localDir, "/") {
		localDir += "/"
	}

	config := servers.Config{
		Host:    "0.0.0.0",
		Port:    11030,
		Context: "",
		DistDir: localDir + "../dist/",
	}

	server, err := servers.New(config)
	if err != nil {
		util.Logger.Error("server new error", zap.Any("config", config), zap.Error(err))
		return
	}

	err = bindFilter(server)
	if err != nil {
		util.Logger.Error("server bind filter error", zap.Any("config", config), zap.Error(err))
		return
	}

	err = bindMapper(server)
	if err != nil {
		util.Logger.Error("server bind mapper error", zap.Any("config", config), zap.Error(err))
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

type Path1Filter struct {
}

func (this_ *Path1Filter) DoFilter(requestContext *servers.HttpRequestContext, chain servers.HttpFilterChain) (err error) {

	util.Logger.Info("path 1 filter start", zap.Any("requestContext", requestContext))
	err = chain.DoFilter(requestContext)
	util.Logger.Info("path 1 filter end", zap.Any("requestContext", requestContext))
	return
}

type Path2Filter struct {
}

func (this_ *Path2Filter) DoFilter(requestContext *servers.HttpRequestContext, chain servers.HttpFilterChain) (err error) {

	util.Logger.Info("path 2 filter start", zap.Any("requestContext", requestContext))
	err = chain.DoFilter(requestContext)
	util.Logger.Info("path 2 filter end", zap.Any("requestContext", requestContext))
	return
}
func bindFilter(server *servers.Server) (err error) {
	err = server.RegisterFilter("/{xx:**}", 1, &Path1Filter{})
	if err != nil {
		return
	}
	err = server.RegisterFilter("/{xx:**}", 2, &Path2Filter{})
	if err != nil {
		return
	}

	return
}
func bindMapper(server *servers.Server) (err error) {
	err = server.RegisterMapper("/data", 0, func(c *servers.HttpRequestContext) (res interface{}, err error) {
		return
	})
	if err != nil {
		return
	}

	return
}
