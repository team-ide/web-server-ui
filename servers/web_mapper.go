package servers

import (
	"reflect"
	"time"
)

type HttpMapperRegister struct {
	mapper            interface{}
	mapperMethod      reflect.Method
	mapperMethodValue reflect.Value
	HttpBaseRegister
}

func (this_ HttpMapperRegister) SetMapper(mapper interface{}) HttpMapperRegister {
	this_.mapper = mapper
	return this_
}

func NewHttpMapperRegister(mapper interface{}, pathPatterns ...string) (register HttpMapperRegister) {
	register = HttpMapperRegister{}
	register.SetMapper(mapper).AddPathPattern(pathPatterns...)
	return
}

func (this_ *Server) processMappers(requestContext *HttpRequestContext) (err error) {
	defer func() {
		requestContext.DoMapperEndTime = time.Now()
		requestContext.PathParams = []*PathParam{}
	}()
	requestContext.DoMapperStartTime = time.Now()

	// 首先判断 是否是静态资源路径 如果是 则直接返回
	isStatic, err := this_.doStatic(requestContext)
	if err != nil {
		return
	}
	if isStatic {
		return
	}

	// 处理 HttpMapper

	pathMatchExtends, err := this_.matchTree(requestContext.Path, this_.mapperPathTree, this_.mapperExcludePathTree)
	if err != nil {
		return
	}

	if len(pathMatchExtends) == 0 {
		this_.doNotFound(requestContext)
		return
	}

	//util.Logger.Info("do mapper match info", zap.Any("path", requestContext.Path), zap.Any("matchList", matchList))
	var callValues []reflect.Value
	for _, one := range pathMatchExtends {
		requestContext.PathParams = one.Params
		mapper := one.Extend.(HttpMapperRegister).mapperMethodValue

		var in []reflect.Value
		callValues = mapper.Call(in)

		var res interface{}
		for _, callValue := range callValues {
			if callValue.Interface() == nil {
				continue
			}
			switch v := callValue.Interface().(type) {
			case error:
				err = v
				break
			default:
				res = v
				break
			}
		}
		if err != nil {
			return
		}
		err = this_.doResult(requestContext, res)
		if err != nil {
			return
		}
	}
	return
}
