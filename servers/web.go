package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
)

func (this_ *Server) bindRouterGroup(routerGroup *gin.RouterGroup) (err error) {
	util.Logger.Info("bind router group start")

	routerGroup.Any("*_requestFullPath", func(c *gin.Context) {
		this_.processRequest(c)
	})

	err = this_.bindStatics()
	if err != nil {
		return
	}
	util.Logger.Info("bind router group end")
	return
}
