package servers

import (
	"fmt"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/web-server-ui/ui"
	"go.uber.org/zap"
	"net"
	"strings"
)

func New(config Config) (ser *Server, err error) {
	ser = &Server{
		config:                     &config,
		filterPathTree:             NewPathTree(""),
		filterExcludePathTree:      NewPathTree(""),
		interceptorPathTree:        NewPathTree(""),
		interceptorExcludePathTree: NewPathTree(""),
		mapperPathTree:             NewPathTree(""),
		mapperExcludePathTree:      NewPathTree(""),
		staticPathCache:            make(map[string]*Static),
		staticNameCache:            make(map[string]*Static),
		pageCache:                  make(map[string]*ui.Page),
	}
	err = ser.init()
	return
}

type Server struct {
	config    *Config
	serverUrl string
	basePath  string

	*ui.App

	webListener net.Listener

	filterPathTree             *PathTree
	filterExcludePathTree      *PathTree
	interceptorPathTree        *PathTree
	interceptorExcludePathTree *PathTree
	mapperPathTree             *PathTree
	mapperExcludePathTree      *PathTree

	staticPathCache map[string]*Static
	staticNameCache map[string]*Static

	pageCache map[string]*ui.Page
}

func (this_ *Server) init() (err error) {
	util.Logger.Info("server init start")
	if this_.config.Host == "" {
		this_.config.Host = "0.0.0.0"
	}

	if this_.config.Port == 0 {
		var listener net.Listener
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			util.Logger.Error("随机端口获取失败", zap.Error(err))
			return
		}
		this_.config.Port = listener.Addr().(*net.TCPAddr).Port
		err = listener.Close()
		if err != nil {
			return
		}
	}
	if this_.config.Context == "" {
		this_.config.Context = "/"
	}
	if !strings.HasPrefix(this_.config.Context, "/") {
		this_.config.Context = "/" + this_.config.Context
	}
	if !strings.HasSuffix(this_.config.Context, "/") {
		this_.config.Context = this_.config.Context + "/"
	}
	if this_.config.DistDir != "" && !strings.HasSuffix(this_.config.DistDir, "/") {
		this_.config.DistDir = this_.config.DistDir + "/"
	}

	if this_.config.Host == "0.0.0.0" || this_.config.Host == ":" || this_.config.Host == "::" {
		this_.serverUrl = fmt.Sprintf("%s://%s:%d", "http", "127.0.0.1", this_.config.Port)
	} else {
		this_.serverUrl = fmt.Sprintf("%s://%s:%d", "http", this_.config.Host, this_.config.Port)
	}
	this_.serverUrl += this_.config.Context
	this_.basePath = this_.config.Context

	if this_.config.Options.SuccessCode == "" {
		this_.config.Options.SuccessCode = "0"
	}

	if this_.config.Options.ErrorCode == "" {
		this_.config.Options.ErrorCode = "-1"
	}

	this_.App = ui.NewApp()
	this_.App.SetBasePath(this_.basePath)

	util.Logger.Info("server info", zap.Any("config", this_.config))
	util.Logger.Info("server info", zap.Any("serverUrl", this_.serverUrl))
	util.Logger.Info("server info", zap.Any("basePath", this_.basePath))
	util.Logger.Info("server init end")
	return
}
