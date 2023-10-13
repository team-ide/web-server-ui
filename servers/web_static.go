package servers

import (
	"bytes"
	"errors"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/web-server-ui/static"
	"github.com/team-ide/web-server-ui/ui"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type StaticPlace int

var (
	staticPlaceStatic StaticPlace = 1
	staticPlaceDir    StaticPlace = 2
	staticPlaceBytes  StaticPlace = 3
	staticPlacePage   StaticPlace = 4
)

type Static struct {
	name  string
	place StaticPlace
	bytes []byte
	page  *ui.Page
}

func (this_ *Static) SetName(name string) *Static {
	this_.name = name
	return this_
}
func (this_ *Static) GetName() string {
	return this_.name
}

func (this_ *Server) doStatic(requestContext *HttpRequestContext) (ok bool, err error) {

	find, ok := this_.staticPathCache[requestContext.Path]
	if !ok {
		return
	}
	err = this_.responseStatic(requestContext, find)
	return
}

func (this_ *Server) responseStatic(requestContext *HttpRequestContext, s *Static) (err error) {

	this_.setHeaderByStaticName(s, requestContext)
	if s.place == staticPlaceDir && strings.HasSuffix(s.name, ".html") {
		var bs []byte
		bs, err = this_.GetOrCopyStatic(s, nil)
		if err != nil {
			return
		}
		this_.writeHtml(requestContext.GetWriter(), bs)
	} else {
		_, err = this_.GetOrCopyStatic(s, requestContext.GetWriter())
	}

	requestContext.Status(http.StatusOK)
	return
}

func (this_ *Server) setHeaderByStaticName(s *Static, requestContext *HttpRequestContext) {

	if s.place == staticPlacePage {
		requestContext.Header("Content-Type", "text/html;charset=UTF-8")
		requestContext.Header("Cache-Control", "no-cache")
	} else if strings.HasSuffix(s.name, ".html") {
		requestContext.Header("Content-Type", "text/html;charset=UTF-8")
		requestContext.Header("Cache-Control", "no-cache")
	} else if strings.HasSuffix(s.name, ".css") {
		requestContext.Header("Content-Type", "text/css;charset=UTF-8")
		// max-age 缓存 过期时间 秒为单位
		requestContext.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(s.name, ".js") {
		requestContext.Header("Content-Type", "application/javascript;charset=UTF-8")
		// max-age 缓存 过期时间 秒为单位
		requestContext.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(s.name, ".woff") ||
		strings.HasSuffix(s.name, ".ttf") ||
		strings.HasSuffix(s.name, ".woff2") ||
		strings.HasSuffix(s.name, ".eot") {
		// max-age 缓存 过期时间 秒为单位
		requestContext.Header("Cache-Control", "max-age=31536000")
	}
}

func (this_ *Server) bindStatics() (err error) {
	util.Logger.Info("bind statics start")

	staticNames := static.GetStaticNames()
	for _, name := range staticNames {
		if name == "index.html" {
			if err = this_.bindStatic("/", "index.html", staticPlaceStatic); err != nil {
				return
			}
		}
		if err = this_.bindStatic("/"+name, name, staticPlaceStatic); err != nil {
			return
		}
	}

	if this_.config.DistDir != "" {
		staticNames, _ = util.LoadDirFilenames(this_.config.DistDir)
		for _, name := range staticNames {
			if name == "index.html" {
				if err = this_.bindStatic("/", "index.html", staticPlaceDir); err != nil {
					return
				}
			}
			if err = this_.bindStatic("/"+name, name, staticPlaceDir); err != nil {
				return
			}
		}
	}

	util.Logger.Info("bind statics end")
	return
}

func (this_ *Server) BindStatic(path, name string, bytes []byte) (err error) {
	s := &Static{
		name:  name,
		place: staticPlaceBytes,
		bytes: bytes,
	}
	this_.staticPathCache[path] = s
	this_.staticNameCache[name] = s
	return
}

func (this_ *Server) BindPage(path, name string, page *ui.Page) (err error) {
	s := &Static{
		name:  name,
		place: staticPlacePage,
		page:  page,
	}
	this_.staticPathCache[path] = s
	this_.staticNameCache[name] = s
	this_.pageCache[name] = page
	return
}

func (this_ *Server) bindStatic(path, name string, place StaticPlace) (err error) {
	s := &Static{
		name:  name,
		place: place,
	}
	this_.staticPathCache[path] = s
	this_.staticNameCache[name] = s
	return
}

func (this_ *Server) GetOrCopyStatic(s *Static, writer io.Writer) (bs []byte, err error) {
	place := s.place
	switch place {
	case staticPlaceDir:
		if writer != nil {
			var f *os.File
			f, err = os.Open(this_.config.DistDir + s.name)
			if err != nil {
				err = errors.New(err.Error() + ",file path:" + this_.config.DistDir + s.name)
				return
			}
			defer func() { _ = f.Close() }()
			_, err = io.Copy(writer, f)
			if err != nil {
				err = errors.New(err.Error() + ",file path:" + this_.config.DistDir + s.name)
				return
			}
			return
		} else {
			bs, err = os.ReadFile(this_.config.DistDir + s.name)
			if err != nil {
				err = errors.New(err.Error() + ",file path:" + this_.config.DistDir + s.name)
				return
			}
			return
		}
	case staticPlacePage:
		if writer != nil {
			options := s.page.NewPageBuilder(writer)
			err = this_.BuildHtml(options)
			return
		} else {
			buffer := &bytes.Buffer{}
			options := s.page.NewPageBuilder(buffer)
			err = this_.BuildHtml(options)
			if err != nil {
				return
			}
			bs = buffer.Bytes()
			return
		}
	case staticPlaceStatic:
		var find bool
		bs, find = static.FindStatic(s.name)
		if !find {
			err = os.ErrNotExist
			return
		}
		break
	case staticPlaceBytes:
		bs = s.bytes
		break
	}
	if writer != nil {
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
