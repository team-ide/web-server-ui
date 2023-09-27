package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"net/http"
	"os"
	"regexp"
	"strings"
	"web-server-ui/static"
)

func (this_ *Server) bindStaticsMapper(routerGroup *gin.RouterGroup) {
	util.Logger.Info("bind statics start")

	this_.bindStaticMapper(routerGroup, "", "index.html")

	if this_.config.DistDir != "" {
		staticNames, _ := util.LoadDirFilenames(this_.config.DistDir)
		for _, name := range staticNames {
			this_.bindStaticMapper(routerGroup, name, name)
		}
	} else {
		staticNames := static.GetStaticNames()
		for _, name := range staticNames {
			this_.bindStaticMapper(routerGroup, name, name)
		}
	}
	util.Logger.Info("bind statics end")
}

func (this_ *Server) bindStaticMapper(routerGroup *gin.RouterGroup, path, name string) {
	util.Logger.Info("bind static", zap.Any("path", path), zap.Any("name", name))
	routerGroup.GET(path, func(c *gin.Context) {
		this_.toStaticByName(c, name)
	})
}

func (this_ *Server) bindStatics(routerGroup *gin.RouterGroup) {
	util.Logger.Info("bind statics start")

	this_.bindStatic(routerGroup, "", "index.html")

	if this_.config.DistDir != "" {
		staticNames, _ := util.LoadDirFilenames(this_.config.DistDir)
		for _, name := range staticNames {
			this_.bindStatic(routerGroup, name, name)
		}
	} else {
		staticNames := static.GetStaticNames()
		for _, name := range staticNames {
			this_.bindStatic(routerGroup, name, name)
		}
	}
	util.Logger.Info("bind statics end")
}

func (this_ *Server) bindStatic(routerGroup *gin.RouterGroup, path, name string) {
	util.Logger.Info("bind static", zap.Any("path", path), zap.Any("name", name))
	routerGroup.GET(path, func(c *gin.Context) {
		this_.toStaticByName(c, name)
	})
}

func (this_ *Server) toStaticByName(c *gin.Context, name string) bool {

	var localFind bool

	var bytes []byte
	if this_.config.DistDir != "" {
		filePath := this_.config.DistDir + name
		localFind, _ = util.PathExists(filePath)
		if localFind {
			bytes, _ = os.ReadFile(filePath)
		}
	}
	if !localFind {
		bytes = static.Asset(name)
		if bytes == nil {
			return false
		}
	}
	this_.setHeaderByName(name, c)
	if strings.HasSuffix(name, ".html") {
		this_.writeHtml(c, name, bytes)
	} else {
		_, _ = c.Writer.Write(bytes)
	}
	c.Status(http.StatusOK)
	return true
}

func (this_ *Server) writeHtml(c *gin.Context, _ string, bytes []byte) {
	templateHTML := string(bytes)

	outHtml := ""
	var re *regexp.Regexp
	re, _ = regexp.Compile(`[$]+{(.+?)}`)
	indexList := re.FindAllIndex(bytes, -1)
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

	_, _ = c.Writer.WriteString(outHtml)
}

func (this_ *Server) setHeaderByName(name string, c *gin.Context) {
	if strings.HasSuffix(name, ".html") {
		c.Header("Content-Type", "text/html")
		c.Header("Cache-Control", "no-cache")
	} else if strings.HasSuffix(name, ".css") {
		c.Header("Content-Type", "text/css")
		// max-age 缓存 过期时间 秒为单位
		c.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(name, ".js") {
		c.Header("Content-Type", "application/javascript")
		// max-age 缓存 过期时间 秒为单位
		c.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(name, ".woff") ||
		strings.HasSuffix(name, ".ttf") ||
		strings.HasSuffix(name, ".woff2") ||
		strings.HasSuffix(name, ".eot") {
		// max-age 缓存 过期时间 秒为单位
		c.Header("Cache-Control", "max-age=31536000")
	}
}
