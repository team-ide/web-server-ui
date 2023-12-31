package servers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
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
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("%s", x))
			util.Logger.Error("process mappers recover error", zap.Any("requestContext", requestContext), zap.Error(err))
		}
		requestContext.DoMapperEndTime = time.Now()
		requestContext.setPathParams(nil)
	}()
	requestContext.DoMapperStartTime = time.Now()

	// 首先判断 是否是静态资源路径 如果是 则直接返回
	isStatic, err := this_.doStatic(requestContext)
	if err != nil {
		util.Logger.Error("process mappers do static error", zap.Any("requestContext", requestContext), zap.Error(err))
		return
	}
	if isStatic {
		return
	}

	// 处理 HttpMapper

	pathMatchExtends, err := this_.matchTree(requestContext.Path, this_.mapperPathTree, this_.mapperExcludePathTree)
	if err != nil {
		util.Logger.Error("process mappers match tree error", zap.Any("requestContext", requestContext), zap.Error(err))
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
		requestContext.setPathParams(one.Params)
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
			util.Logger.Error("process mapper call method error", zap.Any("requestContext", requestContext), zap.Error(err))
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
	var inV interface{}
	for pathVIndex, inType := range inTypes {

		kind := inType.Kind()

		for kind == reflect.Ptr {
			kind = inType.Elem().Kind()
		}
		util.Logger.Debug("get in values", zap.Any("inType", inType), zap.Any("kind", uint(kind)))
		if (kind >= 1 && kind <= 11) || (kind >= 13 && kind <= 14) {
			var pathValue interface{}
			if pathVIndex < requestContext.pathParamValuesSize {
				pathValue = requestContext.pathParamValues[pathVIndex]
			}
			util.Logger.Debug("get in values", zap.Any("pathValue", pathValue))
			inV, err = util.GetValueByType(inType, pathValue)
			if err != nil {
				util.Logger.Error("get in values get value by type error", zap.Any("requestContext", requestContext), zap.Error(err))
				return
			}
		} else {
			name := inType.String()
			switch name {
			case requestContextType1:
				inV = requestContext
				break
			case requestContextType2:
				inV = *requestContext
				break
			default:
				inV, err = util.GetValueByType(inType, "")
				if err != nil {
					util.Logger.Error("get in values get value by type error", zap.Any("requestContext", requestContext), zap.Error(err))
					return
				}
				err = requestContext.c.Bind(&inV)
				if err != nil {
					util.Logger.Error("get in values bind data error", zap.Any("requestContext", requestContext), zap.Error(err))
					return
				}
				break
			}
		}
		inValues = append(inValues, reflect.ValueOf(inV))
	}
	return
}
