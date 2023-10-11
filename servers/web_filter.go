package servers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"time"
)

type HttpFilter interface {
	DoFilter(requestContext *HttpRequestContext, chain HttpFilterChain) (err error)
}

type HttpFilterRegister struct {
	filter HttpFilter
	*HttpBaseRegister
}

func (this_ *HttpFilterRegister) SetFilter(filter HttpFilter) *HttpFilterRegister {
	this_.filter = filter
	return this_
}

func NewHttpFilterRegister(filter HttpFilter, pathPatterns ...string) (register *HttpFilterRegister) {
	register = &HttpFilterRegister{
		HttpBaseRegister: &HttpBaseRegister{},
	}
	register.SetFilter(filter).AddPathPattern(pathPatterns...)
	return
}

type HttpFilterChain interface {
	DoFilter(requestContext *HttpRequestContext) (err error)
}

type HttpFilterChainImpl struct {
	nextFilterIndex int
	filters         []HttpFilter
	filtersSize     int
	pathParamsList  [][]*PathParam
	server          *Server
}

func (this_ *HttpFilterChainImpl) DoFilter(requestContext *HttpRequestContext) (err error) {
	if this_.nextFilterIndex >= this_.filtersSize {
		err = this_.server.processInterceptors(requestContext)
		return
	}
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("%s", x))
		}
	}()
	nextFilter := this_.filters[this_.nextFilterIndex]
	pathParams := this_.pathParamsList[this_.nextFilterIndex]
	this_.nextFilterIndex++
	requestContext.setPathParams(pathParams)
	err = nextFilter.DoFilter(requestContext, this_)
	if err != nil {
		return
	}

	return
}

func (this_ *Server) processFilters(requestContext *HttpRequestContext) (err error) {
	defer func() {
		requestContext.DoFilterEndTime = time.Now()
		requestContext.setPathParams(nil)
	}()
	requestContext.DoFilterStartTime = time.Now()

	var chain = &HttpFilterChainImpl{
		server: this_,
	}
	pathMatchExtends, err := this_.matchTree(requestContext.Path, this_.filterPathTree, this_.filterExcludePathTree)
	if err != nil {
		util.Logger.Error("process filters match tree error", zap.Any("requestContext", requestContext), zap.Error(err))
		return
	}
	//util.Logger.Info("do filter match info", zap.Any("path", requestContext.Path), zap.Any("matchList", matchList))

	var pathParamsList [][]*PathParam
	var filters []HttpFilter

	for _, one := range pathMatchExtends {
		filters = append(filters, one.Extend.(HttpFilterRegister).filter)
		pathParamsList = append(pathParamsList, one.Params)
	}

	// 处理 HttpFilter
	chain.filters = filters
	chain.pathParamsList = pathParamsList
	chain.filtersSize = len(filters)

	err = chain.DoFilter(requestContext)
	if err != nil {
		return
	}

	return
}
