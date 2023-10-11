package servers

import (
	"errors"
	"reflect"
	"strings"
)

type HttpBaseRegister struct {
	pathPatterns        []string
	excludePathPatterns []string
	methods             []string
	order               int
}

func (this_ *HttpBaseRegister) AddPathPattern(pathPatterns ...string) *HttpBaseRegister {
	this_.pathPatterns = append(this_.pathPatterns, pathPatterns...)
	return this_
}

func (this_ *HttpBaseRegister) SetMethod(methods ...string) *HttpBaseRegister {
	this_.methods = append(this_.methods, methods...)
	return this_
}

func (this_ *HttpBaseRegister) ExcludePathPatterns(excludePathPatterns ...string) *HttpBaseRegister {
	this_.excludePathPatterns = append(this_.excludePathPatterns, excludePathPatterns...)
	return this_
}

func (this_ *HttpBaseRegister) SetOrder(order int) *HttpBaseRegister {
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
	objType_ := objType
	for objType_.Kind() == reflect.Ptr {
		objType_ = objType_.Elem()
	}

	fieldNum := objType_.NumField()
	var registers []*HttpMapperRegister
	for i := 0; i < fieldNum; i++ {
		field := objType_.Field(i)
		fieldName := field.Name
		info := GetMapperInfo(field.Tag)
		methodName := info.Bind
		if strings.HasSuffix(fieldName, "Mapper") {
			methodName = fieldName[0 : len(fieldName)-6]
			if methodName == "" {
				err = errors.New("field [" + fieldName + "] mapper method name is empty")
				return
			}
		}
		if methodName == "" {
			continue
		}
		if len(info.Path) == 0 {
			err = errors.New("field [" + fieldName + "] mapper method [" + methodName + "] path empty")
			return
		}

		_, find := objType.MethodByName(methodName)
		if !find {
			err = errors.New("field [" + fieldName + "] mapper method [" + methodName + "] not found")
			return
		}
		methodValue := objValue.MethodByName(methodName)
		register := NewHttpMapperRegister(methodValue.Interface())
		for _, p := range info.Path {
			mapperPath := path + p
			register.AddPathPattern(mapperPath)
		}
		register.SetMethod(info.Method...)

		registers = append(registers, register)
	}
	for _, one := range registers {
		err = this_.RegisterMapper(*one)
		if err != nil {
			return
		}
	}
	return
}

type MapperInfo struct {
	Bind        string   `json:"bind"`
	Method      []string `json:"method"`
	Path        []string `json:"path"`
	ExcludePath []string `json:"excludePath"`
}

func GetMapperInfo(tag reflect.StructTag) *MapperInfo {
	info := &MapperInfo{}

	info.Bind = tag.Get("bind")
	method := tag.Get("method")
	info.Method = strings.Split(method, ",")
	path := tag.Get("path")
	info.Path = strings.Split(path, ",")
	excludePath := tag.Get("excludePath")
	info.ExcludePath = strings.Split(excludePath, ",")

	return info
}
