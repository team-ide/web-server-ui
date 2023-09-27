package servers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"strings"
)

func (this_ *Server) bindRouterGroup(routerGroup *gin.RouterGroup) {
	util.Logger.Info("bind router group start")

	this_.bindStatics(routerGroup)

	routerGroup.Any("*_requestFullPath", func(c *gin.Context) {
		_requestFullPath := c.Params.ByName("_requestFullPath")
		fmt.Println("_requestFullPath:", _requestFullPath)

		// 处理 HttpFilter

		// 处理 HttpHandlerInterceptor

		// 处理 HttpMapper

		mapperMatch, err := this_.mapperPathTree.Match(_requestFullPath)
		if err != nil {
			return
		}
		if mapperMatch.Matched {
			httpMapper := mapperMatch.Node.GetExtend().(HttpMapper)
			httpMapper()
		}
	})
	util.Logger.Info("bind router group end")
}

type HttpRequest struct {
}

type HttpResponse struct {
}

type HttpMapper func(request *HttpRequest)

type HttpFilter interface {
	DoFilter(request *HttpRequest, response *HttpResponse, chain *HttpFilterChain) (err error)
}

type HttpFilterChain struct {
	nextFilterIndex int
	filters         []HttpFilter
}

func (this_ *HttpFilterChain) DoFilter(request *HttpRequest, response *HttpResponse) (err error) {

	return
}

type HttpHandlerInterceptor interface {
	PreHandle(request *HttpRequest, response *HttpResponse) bool
	PostHandle(request *HttpRequest, response *HttpResponse)
	AfterCompletion(request *HttpRequest, response *HttpResponse)
}

func (this_ *Server) RegisterHttpFilter(matchPath string, order int, filter HttpFilter) (err error) {

	matchPaths, err := this_.validateMatchPath(matchPath)
	if err != nil {
		util.Logger.Error("validateMatchPath error", zap.Error(err))
		return
	}

	for _, path := range matchPaths {
		err = this_.filterPathTree.AddPath(path, order, filter)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *Server) RegisterMapper(matchPath string, order int, mapper HttpMapper) (err error) {

	matchPaths, err := this_.validateMatchPath(matchPath)
	if err != nil {
		util.Logger.Error("validateMatchPath error", zap.Error(err))
		return
	}

	for _, path := range matchPaths {
		err = this_.mapperPathTree.AddPath(path, order, mapper)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *Server) RegisterHttpHandlerInterceptor(matchPath string, order int, interceptor HttpHandlerInterceptor) (err error) {

	matchPaths, err := this_.validateMatchPath(matchPath)
	if err != nil {
		util.Logger.Error("validateMatchPath error", zap.Error(err))
		return
	}
	for _, path := range matchPaths {
		err = this_.handlerInterceptorPathTree.AddPath(path, order, interceptor)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *Server) validateMatchPath(matchPath string) (matchPaths []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%s", e))
		}
	}()
	if matchPath == "" {
		err = errors.New("match path is empty")
		return
	}
	ss := strings.Split(matchPath, ",")

	for _, path := range ss {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		matchPaths = append(matchPaths, path)
	}
	if len(matchPaths) == 0 {
		err = errors.New("match path is empty")
		return
	}

	return
}
