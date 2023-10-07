package servers

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"os"
	"regexp"
	"strings"
	"web-server-ui/static"
)

func (this_ *Server) bindStaticsMapper() (err error) {
	util.Logger.Info("bind statics start")

	if err = this_.bindStaticMapper("/", "index.html"); err != nil {
		return
	}

	if this_.config.DistDir != "" {
		staticNames, _ := util.LoadDirFilenames(this_.config.DistDir)
		for _, name := range staticNames {
			if err = this_.bindStaticMapper("/"+name, name); err != nil {
				return
			}
		}
	} else {
		staticNames := static.GetStaticNames()
		for _, name := range staticNames {
			if err = this_.bindStaticMapper("/"+name, name); err != nil {
				return
			}
		}
	}
	util.Logger.Info("bind statics end")
	return
}

func (this_ *Server) bindStaticMapper(path, name string) (err error) {
	util.Logger.Info("bind static", zap.Any("path", path), zap.Any("name", name))
	err = this_.RegisterMapper(path, 0, func(c *HttpRequestContext) (res interface{}, err error) {
		res = &ResultPage{
			Page: name,
		}
		return
	})
	if err != nil {
		util.Logger.Error("bing static mapper error", zap.Error(err))
		return
	}
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
