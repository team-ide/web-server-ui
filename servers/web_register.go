package servers

import (
	"errors"
	"fmt"
	"reflect"
)

type HttpBaseRegister struct {
	pathPatterns        []string
	excludePathPatterns []string
	methods             []string
	order               int
}

func (this_ HttpBaseRegister) AddPathPattern(pathPatterns ...string) HttpBaseRegister {
	this_.pathPatterns = append(this_.pathPatterns, pathPatterns...)
	return this_
}

func (this_ HttpBaseRegister) SetMethod(methods ...string) HttpBaseRegister {
	this_.methods = append(this_.methods, methods...)
	return this_
}

func (this_ HttpBaseRegister) ExcludePathPatterns(excludePathPatterns ...string) HttpBaseRegister {
	this_.excludePathPatterns = append(this_.excludePathPatterns, excludePathPatterns...)
	return this_
}

func (this_ HttpBaseRegister) SetOrder(order int) HttpBaseRegister {
	this_.order = order
	return this_
}

// RegisterFilter
// 顺序是 order 从小到大 执行
func (this_ *Server) RegisterFilter(register HttpFilterRegister) (err error) {
	if len(register.pathPatterns) == 0 {
		err = errors.New("register path patterns is nil")
		return
	}
	if register.filter == nil {
		err = errors.New("register filter is nil")
		return
	}

	for _, path := range register.pathPatterns {
		err = this_.filterPathTree.AddPath(path, register.order, register)
		if err != nil {
			return
		}
	}
	for _, path := range register.excludePathPatterns {
		err = this_.filterExcludePathTree.AddPath(path, register.order, register)
		if err != nil {
			return
		}
	}

	return
}

// RegisterInterceptor
// 顺序是 order 从小到大 执行
func (this_ *Server) RegisterInterceptor(register HttpInterceptorRegister) (err error) {
	if len(register.pathPatterns) == 0 {
		err = errors.New("register path patterns is nil")
		return
	}
	if register.interceptor == nil {
		err = errors.New("register interceptor is nil")
		return
	}

	for _, path := range register.pathPatterns {
		err = this_.interceptorPathTree.AddPath(path, register.order, register)
		if err != nil {
			return
		}
	}
	for _, path := range register.excludePathPatterns {
		err = this_.interceptorExcludePathTree.AddPath(path, register.order, register)
		if err != nil {
			return
		}
	}

	return
}

// RegisterMapper
// 顺序是 order 从小到大 执行
func (this_ *Server) RegisterMapper(register HttpMapperRegister) (err error) {
	if len(register.pathPatterns) == 0 {
		err = errors.New("register path patterns is nil")
		return
	}
	if register.mapper == nil {
		err = errors.New("register mapper is nil")
		return
	}

	for _, path := range register.pathPatterns {
		err = this_.mapperPathTree.AddPath(path, register.order, register)
		if err != nil {
			return
		}
	}
	for _, path := range register.excludePathPatterns {
		err = this_.mapperExcludePathTree.AddPath(path, register.order, register)
		if err != nil {
			return
		}
	}

	return
}

// RegisterMapperObj
// 顺序是 order 从小到大 执行
func (this_ *Server) RegisterMapperObj(path string, mapperObj interface{}) (err error) {
	objType := reflect.TypeOf(mapperObj)
	objValue := reflect.ValueOf(mapperObj)

	methodNum := objType.NumMethod()
	for i := 0; i < methodNum; i++ {
		methodType := objType.Method(i)
		methodValue := objValue.Method(i)
		fmt.Println("method:", methodType.Name, ",doc:", methodValue)
		mapperPath := path + methodType.Name
		register := NewHttpMapperRegister(methodValue.Interface())
		register.AddPathPattern(mapperPath)
		err = this_.RegisterMapper(register)
		if err != nil {
			return
		}
	}
	return
}
