package servers

import "time"

type HttpInterceptor interface {
	Before(requestContext *HttpRequestContext) bool
	After(requestContext *HttpRequestContext)
}

type HttpInterceptorRegister struct {
	interceptor HttpInterceptor
	HttpBaseRegister
}

func (this_ HttpInterceptorRegister) SetInterceptor(interceptor HttpInterceptor) HttpInterceptorRegister {
	this_.interceptor = interceptor
	return this_
}

func NewHttpInterceptorRegister(interceptor HttpInterceptor, pathPatterns ...string) (register HttpInterceptorRegister) {
	register = HttpInterceptorRegister{}
	register.SetInterceptor(interceptor).AddPathPattern(pathPatterns...)
	return
}

func (this_ *Server) processInterceptors(requestContext *HttpRequestContext) (err error) {
	defer func() {
		requestContext.DoInterceptorEndTime = time.Now()
		requestContext.PathParams = []*PathParam{}
	}()
	requestContext.DoInterceptorStartTime = time.Now()

	// 处理 HttpHandlerInterceptor

	pathMatchExtends, err := this_.matchTree(requestContext.Path, this_.interceptorPathTree, this_.interceptorExcludePathTree)
	if err != nil {
		return
	}
	//util.Logger.Info("do interceptor match info", zap.Any("path", requestContext.Path), zap.Any("matchList", matchList))
	for _, one := range pathMatchExtends {
		requestContext.PathParams = one.Params
		interceptor := one.Extend.(HttpInterceptorRegister).interceptor
		if !interceptor.Before(requestContext) {
			return
		}
	}
	err = this_.processMappers(requestContext)

	for _, one := range pathMatchExtends {
		requestContext.PathParams = one.Params
		interceptor := one.Extend.(HttpInterceptorRegister).interceptor
		interceptor.After(requestContext)
	}
	return
}
