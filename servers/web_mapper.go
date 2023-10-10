package servers

import (
	"reflect"
	"time"
)

type HttpMapperRegister struct {
	mapper      interface{}
	mapperType  reflect.Type
	mapperValue reflect.Value
	inTypes     []reflect.Type
	outTypes    []reflect.Type
	*HttpBaseRegister
}

func (this_ *HttpMapperRegister) setMapper(mapper interface{}) *HttpMapperRegister {
	this_.mapper = mapper

	this_.mapperType = reflect.TypeOf(mapper)
	this_.mapperValue = reflect.ValueOf(mapper)
	numIn := this_.mapperType.NumIn()
	numOut := this_.mapperType.NumOut()
	for i := 0; i < numIn; i++ {
		inType := this_.mapperType.In(i)
		this_.inTypes = append(this_.inTypes, inType)
	}
	for i := 0; i < numOut; i++ {
		outType := this_.mapperType.Out(i)
		this_.outTypes = append(this_.outTypes, outType)
	}
	//fmt.Println("mapper:", mapperType, ",numIn:", numIn, ",numOut:", numOut)
	return this_
}

func NewHttpMapperRegister(mapper interface{}, pathPatterns ...string) (register *HttpMapperRegister) {
	register = &HttpMapperRegister{
		HttpBaseRegister: &HttpBaseRegister{},
	}
	register.setMapper(mapper).AddPathPattern(pathPatterns...)

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
	var inValues []reflect.Value
	for _, one := range pathMatchExtends {
		requestContext.PathParams = one.Params
		mapperRegister := one.Extend.(HttpMapperRegister)

		inValues, err = this_.GetInValues(requestContext, mapperRegister.inTypes)
		callValues = mapperRegister.mapperValue.Call(inValues)

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

var (
	requestContextType1 = reflect.TypeOf(&HttpRequestContext{}).String()
	requestContextType2 = reflect.TypeOf(HttpRequestContext{}).String()
)

func (this_ *Server) GetInValues(requestContext *HttpRequestContext, inTypes []reflect.Type) (inValues []reflect.Value, err error) {
	for _, inType := range inTypes {

		var inValue reflect.Value
		name := inType.String()
		switch name {
		case requestContextType1:
			inValue = reflect.ValueOf(requestContext)
			break
		case requestContextType2:
			inValue = reflect.ValueOf(*requestContext)
			break
		}

		inValues = append(inValues, inValue)
	}
	return
}
