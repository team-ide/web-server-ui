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
	reg := servers.NewHttpFilterRegister(&Path1Filter{})
	reg.AddPathPattern("/{xx:**}").SetOrder(1)
	err = server.RegisterFilter(*reg)
	if err != nil {
		return
	}

	reg = servers.NewHttpFilterRegister(&Path2Filter{})
	reg.AddPathPattern("/{xx:**}").SetOrder(2)
	err = server.RegisterFilter(*reg)
	if err != nil {
		return
	}

	return
}
func bindMapper(server *servers.Server) (err error) {
	//err = server.RegisterMapper("/data", 0, func(requestContext *servers.HttpRequestContext) (res interface{}, err error) {
	//
	//	res = servers.NewResultData("ok")
	//	return
	//})
	//if err != nil {
	//	return
	//}

	err = server.RegisterMapperObj("/user", &UserMapper{
		server: server,
	})
	if err != nil {
		return
	}

	return
}

type UserMapper struct {
	server       *servers.Server
	IndexMapper  string `path:"/index" method:"get"`
	GetMapper    string `path:"/get/{userId}" method:"get"`
	InsertMapper string `path:"/insert" method:"post"`
	UpdateMapper string `path:"/update/{userId}" method:"post"`
	DeleteMapper string `path:"/delete/{userId}" method:"post"`
}

// Index
// mapper:/index
func (this_ *UserMapper) Index(requestContext *servers.HttpRequestContext) (res interface{}, err error) {

	page := this_.server.NewPage()
	res = this_.server.NewResultPage(page)
	return
}

// Get
// mapper:/get/{userId}
func (this_ *UserMapper) Get(userId int64, requestContext *servers.HttpRequestContext) (res interface{}, err error) {
	fmt.Println("Get userId:", userId)
	res = userId
	return
}

// Insert
// mapper:/insert
func (this_ *UserMapper) Insert(requestContext *servers.HttpRequestContext, userInfo *UserInfo) (res interface{}, err error) {
	fmt.Println("Insert userInfo:", util.GetStringValue(userInfo))
	res = userInfo
	return
}

// Update
// mapper:/update/{userId}
func (this_ *UserMapper) Update(userId int64, requestContext *servers.HttpRequestContext, userInfo *UserInfo) (res interface{}, err error) {
	fmt.Println("Update userId:", userId)
	fmt.Println("Update userInfo:", util.GetStringValue(userInfo))

	return
}

// Delete
// mapper:/delete/{userId}
func (this_ *UserMapper) Delete(userId int64, requestContext *servers.HttpRequestContext) (res interface{}, err error) {
	fmt.Println("Delete userId:", userId)
	return
}

type UserInfo struct {
	Name string `json:"name"`
	Age  uint   `json:"age"`
}
