package servers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"regexp"
	"sort"
	"strings"
)

var (
	// 去除多余斜杠
	trimMoreSlash, _ = regexp.Compile("/+")
)

func (this_ *Server) bindRouterGroup(routerGroup *gin.RouterGroup) {
	util.Logger.Info("bind router group start")

	this_.bindStatics(routerGroup)

	routerGroup.Any("*_requestFullPath", func(c *gin.Context) {
		_requestFullPath := c.Params.ByName("_requestFullPath")
		fmt.Println("_requestFullPath:", _requestFullPath)
		requestPath := trimMoreSlash.ReplaceAllLiteralString(_requestFullPath, "/")
		fmt.Println("requestPath:", requestPath)

	})
	util.Logger.Info("bind router group end")
}

type HttpRequest struct {
}

type HttpResponse struct {
}

type HttpFilter interface {
	DoFilter(request *HttpRequest, response *HttpResponse) (err error)
}

type HttpHandlerInterceptor interface {
	PreHandle(request *HttpRequest, response *HttpResponse) bool
	PostHandle(request *HttpRequest, response *HttpResponse)
	AfterCompletion(request *HttpRequest, response *HttpResponse)
}

type RegisterHttpFilter struct {
	matchPaths   []string
	matchRegexps []*regexp.Regexp
	order        int
	filter       HttpFilter
}

type RegisterHttpHandlerInterceptor struct {
	matchPaths   []string
	matchRegexps []*regexp.Regexp
	order        int
	interceptor  HttpHandlerInterceptor
}

func (this_ *Server) RegisterHttpFilter(matchPath string, order int, filter HttpFilter) (err error) {

	matchPaths, matchRegexps, err := this_.validateMatchPath(matchPath)
	if err != nil {
		util.Logger.Error("validateMatchPath error", zap.Error(err))
		return
	}

	register := &RegisterHttpFilter{
		matchPaths:   matchPaths,
		matchRegexps: matchRegexps,
		order:        order,
		filter:       filter,
	}
	this_.registerHttpFilters = append(this_.registerHttpFilters, register)

	// Order 正序
	sort.Slice(this_.registerHttpFilters, func(i, j int) bool {
		return this_.registerHttpFilters[i].order < this_.registerHttpFilters[j].order
	})

	return
}

func (this_ *Server) RegisterHttpHandlerInterceptor(matchPath string, order int, interceptor HttpHandlerInterceptor) (err error) {

	matchPaths, matchRegexps, err := this_.validateMatchPath(matchPath)
	if err != nil {
		util.Logger.Error("validateMatchPath error", zap.Error(err))
		return
	}

	register := &RegisterHttpHandlerInterceptor{
		matchPaths:   matchPaths,
		matchRegexps: matchRegexps,
		order:        order,
		interceptor:  interceptor,
	}
	this_.registerHttpHandlerInterceptors = append(this_.registerHttpHandlerInterceptors, register)

	// Order 正序
	sort.Slice(this_.registerHttpFilters, func(i, j int) bool {
		return this_.registerHttpFilters[i].order < this_.registerHttpFilters[j].order
	})

	return
}

func (this_ *Server) validateMatchPath(matchPath string) (matchPaths []string, matchRegexps []*regexp.Regexp, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%s", e))
		}
	}()
	if matchPath == "" {
		err = errors.New("matchPath is empty")
		return
	}
	ss := strings.Split(matchPath, ",")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		reg := regexp.MustCompile(s)
		if reg == nil {
			err = errors.New("path [" + s + "] regexp MustCompile error")
			return
		}
		matchPaths = append(matchPaths, s)
		matchRegexps = append(matchRegexps, reg)
	}
	if len(matchPaths) == 0 {
		err = errors.New("matchPath [" + matchPath + "] not to regexps")
		return
	}

	return
}
