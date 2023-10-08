package servers

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"web-server-ui/static"
)

type Static struct {
	Name string `json:"name"`
}

func (this_ *Server) doStatic(requestContext *HttpRequestContext) (ok bool, err error) {

	find, ok := this_.staticPathCache[requestContext.Path]
	if !ok {
		return
	}

	util.Logger.Info("do static", zap.Any("name", find.Name))
	this_.setHeaderByStaticName(find.Name, requestContext)
	if strings.HasSuffix(find.Name, ".html") {
		var bs []byte
		bs, err = this_.ReadStatic(find.Name)
		if err != nil {
			return
		}
		this_.writeHtml(requestContext.GetWriter(), bs)
	} else {
		err = this_.CopyStatic(find.Name, requestContext.GetWriter())
	}

	requestContext.Status(http.StatusOK)
	return
}

func (this_ *Server) setHeaderByStaticName(name string, requestContext *HttpRequestContext) {
	if strings.HasSuffix(name, ".html") {
		requestContext.Header("Content-Type", "text/html")
		requestContext.Header("Cache-Control", "no-cache")
	} else if strings.HasSuffix(name, ".css") {
		requestContext.Header("Content-Type", "text/css")
		// max-age 缓存 过期时间 秒为单位
		requestContext.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(name, ".js") {
		requestContext.Header("Content-Type", "application/javascript")
		// max-age 缓存 过期时间 秒为单位
		requestContext.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(name, ".woff") ||
		strings.HasSuffix(name, ".ttf") ||
		strings.HasSuffix(name, ".woff2") ||
		strings.HasSuffix(name, ".eot") {
		// max-age 缓存 过期时间 秒为单位
		requestContext.Header("Cache-Control", "max-age=31536000")
	}
}

func (this_ *Server) bindStatics() (err error) {
	util.Logger.Info("bind statics start")

	if err = this_.bindStatic("/", "index.html"); err != nil {
		return
	}

	if this_.config.DistDir != "" {
		staticNames, _ := util.LoadDirFilenames(this_.config.DistDir)
		for _, name := range staticNames {
			if err = this_.bindStatic("/"+name, name); err != nil {
				return
			}
		}
	} else {
		staticNames := static.GetStaticNames()
		for _, name := range staticNames {
			if err = this_.bindStatic("/"+name, name); err != nil {
				return
			}
		}
	}
	util.Logger.Info("bind statics end")
	return
}

func (this_ *Server) bindStatic(path, name string) (err error) {
	s := &Static{
		Name: name,
	}
	this_.staticPathCache[path] = s
	return
}

func (this_ *Server) ReadStatic(path string) (bs []byte, err error) {
	if this_.config.DistDir != "" {
		bs, err = os.ReadFile(this_.config.DistDir + path)
		if err != nil {
			return
		}
	} else {
		var find bool
		bs, find = static.FindStatic(path)
		if !find {
			err = os.ErrNotExist
			return
		}
	}
	return
}

func (this_ *Server) CopyStatic(path string, writer io.Writer) (err error) {
	if this_.config.DistDir != "" {
		var f *os.File
		f, err = os.Open(this_.config.DistDir + path)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()

		_, err = io.Copy(writer, f)

	} else {
		bs, find := static.FindStatic(path)
		if !find {
			err = os.ErrNotExist
			return
		}
		_, err = writer.Write(bs)
	}
	return
}

func (this_ *Server) writeHtml(writer io.Writer, bs []byte) {
	templateHTML := string(bs)

	outHtml := ""
	var re *regexp.Regexp
	re, _ = regexp.Compile(`[$]+{(.+?)}`)
	indexList := re.FindAllIndex(bs, -1)
	var lastIndex int = 0
	for _, indexes := range indexList {
		outHtml += templateHTML[lastIndex:indexes[0]]

		lastIndex = indexes[1]

		script := strings.TrimSpace(templateHTML[indexes[0]+2 : indexes[1]-1])

		var scriptValue string

		switch script {
		case "":
			break
		case "basePath":
			scriptValue = this_.basePath
			break
		}
		outHtml += scriptValue
	}
	outHtml += templateHTML[lastIndex:]

	_, _ = writer.Write([]byte(outHtml))
}
