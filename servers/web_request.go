package servers

import (
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

	pathParams          []*PathParam
	pathParamMap        map[string]string
	pathParamValues     []string
	pathParamValuesSize int
}

func (this_ *HttpRequestContext) setPathParams(pathParams []*PathParam) *HttpRequestContext {
	this_.pathParamMap = make(map[string]string)
	this_.pathParamValues = []string{}
	this_.pathParamValuesSize = len(pathParams)
	if this_.pathParamValuesSize == 0 {
		return this_
	}
	for _, one := range pathParams {
		this_.pathParamValues = append(this_.pathParamValues, one.Value)
		_, f := this_.pathParamMap[one.Name]
		if !f {
			//fmt.Println("path param name:", one.Name, ",value:", one.Value)
			this_.pathParamMap[one.Name] = one.Value
		}
	}
	return this_
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

func (this_ *Server) processRequest(c *gin.Context) {
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
			util.Logger.Error("process request recover error", zap.Any("requestContext", requestContext), zap.Error(err))
		}
		requestContext.EndTime = time.Now()

		if err != nil {
			this_.doError(requestContext, err)
			return
		}
	}()

	err = this_.processFilters(requestContext)

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
		util.Logger.Error("process mapper do result static error", zap.Any("requestContext", requestContext), zap.Error(err))
		return
	}

	return
}

func (this_ *Server) doResultData(requestContext *HttpRequestContext, data *ResultData) (err error) {
	if data.Code == "" {
		data.Code = this_.config.Options.SuccessCode
	}
	requestContext.c.JSON(http.StatusOK, data)
	return
}

func (this_ *Server) doNotFound(requestContext *HttpRequestContext) {
	util.Logger.Warn("http request not found", zap.Any("requestContext", requestContext))

	requestContext.Status(http.StatusNotFound)
	return
}

func (this_ *Server) doError(requestContext *HttpRequestContext, err error) {
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
