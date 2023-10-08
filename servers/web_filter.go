package servers

import (
	"errors"
	"fmt"
)

type HttpFilter interface {
	DoFilter(requestContext *HttpRequestContext, chain HttpFilterChain) (err error)
}

type HttpFilterChain interface {
	DoFilter(requestContext *HttpRequestContext) (err error)
}

type HttpFilterChainImpl struct {
	nextFilterIndex int
	filters         []HttpFilter
	filtersSize     int
	server          *Server
}

func (this_ *HttpFilterChainImpl) DoFilter(requestContext *HttpRequestContext) (err error) {
	if this_.nextFilterIndex >= this_.filtersSize {
		err = this_.server.doInterceptor(requestContext)
		return
	}
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("%s", x))
		}
	}()
	nextFilter := this_.filters[this_.nextFilterIndex]
	this_.nextFilterIndex++
	err = nextFilter.DoFilter(requestContext, this_)
	if err != nil {
		return
	}

	return
}
