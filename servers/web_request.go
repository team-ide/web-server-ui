package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type HttpRequestContext struct {
	Path      string    `json:"path,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`

	DoFilterStartTime time.Time `json:"doFilterStartTime,omitempty"`
	DoFilterEndTime   time.Time `json:"doFilterEndTime,omitempty"`

	DoInterceptorStartTime time.Time `json:"doInterceptorStartTime,omitempty"`
	DoInterceptorEndTime   time.Time `json:"doInterceptorEndTime,omitempty"`

	DoMapperStartTime time.Time `json:"doMapperStartTime,omitempty"`
	DoMapperEndTime   time.Time `json:"doMapperEndTime,omitempty"`
	c                 *gin.Context

	PathParams []*PathParam `json:"pathParams,omitempty"`
}

func (this_ *HttpRequestContext) Write(bs []byte) (int, error) {
	return this_.c.Writer.Write(bs)
}

func (this_ *HttpRequestContext) WriteString(str string) (int, error) {
	return this_.c.Writer.WriteString(str)
}

func (this_ *HttpRequestContext) Header(key, value string) {
	this_.c.Header(key, value)
}

func (this_ *HttpRequestContext) Status(status int) {
	this_.c.Status(status)
}

func (this_ *HttpRequestContext) GetWriter() io.Writer {
	return this_.c.Writer
}

func (this_ *Server) doRequest(c *gin.Context) {
	_requestFullPath := c.Params.ByName("_requestFullPath")
	var err error

	var requestContext = &HttpRequestContext{
		Path:      _requestFullPath,
		StartTime: time.Now(),
		c:         c,
	}

	defer func() {
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("%s", x))
		}
		requestContext.EndTime = time.Now()

		if err != nil {
			this_.doError(requestContext, err)
			return
		}
	}()

	err = this_.doFilter(requestContext)

	return
}

func (this_ *Server) doFilter(requestContext *HttpRequestContext) (err error) {
	defer func() {
		requestContext.DoFilterEndTime = time.Now()
		requestContext.PathParams = []*PathParam{}
	}()
	requestContext.DoFilterStartTime = time.Now()

	var chain = &HttpFilterChainImpl{
		server: this_,
	}
	var filters []HttpFilter
	matchList, err := this_.filterPathTree.Match(requestContext.Path)
	if err != nil {
		return
	}
	//util.Logger.Info("do filter match info", zap.Any("path", requestContext.Path), zap.Any("matchList", matchList))

	var pathParamsList [][]*PathParam

	for _, one := range matchList {

		es := one.Node.GetExtends()
		requestContext.PathParams = one.Params
		for _, e := range es {
			filters = append(filters, e.GetExtend().(HttpFilter))
			pathParamsList = append(pathParamsList, one.Params)
		}
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

func (this_ *Server) doInterceptor(requestContext *HttpRequestContext) (err error) {
	defer func() {
		requestContext.DoInterceptorEndTime = time.Now()
		requestContext.PathParams = []*PathParam{}
	}()
	requestContext.DoInterceptorStartTime = time.Now()

	// 处理 HttpHandlerInterceptor

	matchList, err := this_.interceptorPathTree.Match(requestContext.Path)
	if err != nil {
		return
	}
	//util.Logger.Info("do interceptor match info", zap.Any("path", requestContext.Path), zap.Any("matchList", matchList))
	for _, one := range matchList {

		es := one.Node.GetExtends()
		requestContext.PathParams = one.Params
		for _, e := range es {
			interceptor := e.GetExtend().(HttpInterceptor)
			if !interceptor.PreHandle(requestContext) {
				return
			}
		}
	}

	err = this_.doMapper(requestContext)
	return
}

func (this_ *Server) doMapper(requestContext *HttpRequestContext) (err error) {
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

	matchList, err := this_.mapperPathTree.Match(requestContext.Path)
	if err != nil {
		return
	}

	if len(matchList) == 0 {
		this_.doNotFound(requestContext)
		return
	}

	//util.Logger.Info("do mapper match info", zap.Any("path", requestContext.Path), zap.Any("matchList", matchList))
	var res interface{}
	for _, one := range matchList {
		es := one.Node.GetExtends()
		requestContext.PathParams = one.Params
		for _, e := range es {
			mapper := e.GetExtend().(HttpMapper)
			res, err = mapper(requestContext)
			if err != nil {
				return
			}
			err = this_.doResult(requestContext, res)
			if err != nil {
				return
			}
		}
	}
	return
}

func (this_ *Server) doResult(requestContext *HttpRequestContext, result interface{}) (err error) {
	if result == ResultNone {
		return
	}
	switch t := result.(type) {
	case ResultStatic:
		err = this_.doResultStatic(requestContext, &t)
		break
	case *ResultStatic:
		err = this_.doResultStatic(requestContext, t)
		break
	case ResultData:
		err = this_.doResultData(requestContext, &t)
		break
	case *ResultData:
		err = this_.doResultData(requestContext, t)
		break
	default:
		err = this_.doResultData(requestContext, &ResultData{
			Data: result,
		})
		break
	}
	requestContext.Status(http.StatusOK)
	return
}

func (this_ *Server) doResultStatic(requestContext *HttpRequestContext, s *ResultStatic) (err error) {

	err = this_.responseStatic(requestContext, s.Name)
	if err != nil {
		return
	}

	return
}

func (this_ *Server) doResultData(requestContext *HttpRequestContext, data *ResultData) (err error) {
	if data.Code == "" {
		data.Code = this_.config.Options.SuccessCode
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, err = requestContext.Write(bs)

	return
}

func (this_ *Server) doNotFound(requestContext *HttpRequestContext) {
	util.Logger.Warn("http request not found", zap.Any("requestContext", requestContext))

	requestContext.Status(http.StatusNotFound)
	return
}

func (this_ *Server) doError(requestContext *HttpRequestContext, err error) {
	util.Logger.Error("http request error", zap.Any("requestContext", requestContext), zap.Error(err))

	var code = this_.config.Options.ErrorCode

	var a *CodeError
	if errors.As(err, &a) {
		code = a.Code
	}

	data := &ResultData{
		Code: code,
		Msg:  "request path [" + requestContext.Path + "] error," + err.Error(),
	}

	_ = this_.doResultData(requestContext, data)

	return
}

type ResultStatic struct {
	Name string `json:"name"`
}

type ResultData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type CodeError struct {
	Code string `json:"code"`
	Err  error  `json:"err"`
}

func (this_ *CodeError) Error() string {
	var err string
	if this_.Err != nil {
		err = this_.Err.Error()
	}
	return err
}

// NewResultStatic 返回静态资源
func NewResultStatic(name string) *ResultStatic {
	return &ResultStatic{
		Name: name,
	}
}

// NewResultError 返回错误信息
func NewResultError(code string, err error) *ResultData {
	return &ResultData{
		Code: code,
		Msg:  err.Error(),
	}
}

// NewResultData 返回结果
func NewResultData(data interface{}) *ResultData {
	return &ResultData{
		Data: data,
	}
}

// NewCodeError 返回带错误码的异常
func NewCodeError(code string, err error) error {
	return &CodeError{
		Code: code,
		Err:  err,
	}
}

var (
	// ResultNone 返回该结果不会通过response写入任何数据
	ResultNone = new(int)
)
