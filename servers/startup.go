package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

func (this_ *Server) Startup(onDown func()) (err error) {
	util.Logger.Info("server startup start")
	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	router.MaxMultipartMemory = (1024 * 50) << 20 // 设置最大上传大小为50G

	routerGroup := router.Group(this_.config.Context)

	this_.bindStatics(routerGroup)

	var ins []net.Interface
	ins, err = net.Interfaces()
	if err != nil {
		return
	}

	out := ""
	out += fmt.Sprintf("服务启动，访问地址:\n")
	if this_.config.Host == "0.0.0.0" || this_.config.Host == ":" || this_.config.Host == "::" {
		out += fmt.Sprintf("\t%s://%s:%d%s\n", "http", "127.0.0.1", this_.config.Port, this_.config.Context)
		for _, in := range ins {
			if in.Flags&net.FlagUp == 0 {
				continue
			}
			if in.Flags&net.FlagLoopback != 0 {
				continue
			}
			var adders []net.Addr
			adders, err = in.Addrs()
			if err != nil {
				return
			}
			for _, addr := range adders {
				ip := util.GetIpFromAddr(addr)
				if ip == nil {
					continue
				}
				out += fmt.Sprintf("\t%s://%s:%d%s\n", "http", ip.String(), this_.config.Port, this_.config.Context)
			}
		}
	} else {
		out += fmt.Sprintf("\t%s://%s:%d%s\n", "http", this_.config.Host, this_.config.Port, this_.config.Context)
	}

	fmt.Println(out)
	addr := fmt.Sprintf("%s:%d", this_.config.Host, this_.config.Port)
	util.Logger.Info("http server start", zap.Any("addr", addr))
	s := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	this_.webListener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return
	}
	go func() {
		err = s.Serve(this_.webListener)
		if err != nil {
			util.Logger.Error("Web启动失败", zap.Error(err))
		}
		if onDown != nil {
			onDown()
		}
	}()
	var checkStartTime = util.GetNowMilli()
	for {
		var newTime = util.GetNowMilli()
		if (newTime - checkStartTime) > 1000*5 {
			util.Logger.Warn("服务启动检查超过5秒，不再检测")
			break
		}
		time.Sleep(time.Millisecond * 100)
		checkURL := this_.serverUrl
		util.Logger.Info("监听服务是否启动成功", zap.Any("checkURL", checkURL))
		res, e := http.Get(checkURL)
		if e != nil {
			util.Logger.Warn("监听服务连接失败，将继续监听", zap.Any("checkURL", checkURL), zap.Any("error", e.Error()))
			continue
		}
		if res == nil {
			util.Logger.Warn("监听服务连接无返回，不再监听", zap.Any("checkURL", checkURL))
			continue
		}
		if res.StatusCode == 200 {
			_ = res.Body.Close()
			util.Logger.Info("服务启动成功", zap.Any("serverUrl", this_.serverUrl))
			break
		}
		util.Logger.Info("服务未启动完成", zap.Any("statusCode", res.StatusCode))
	}
	util.Logger.Info("server startup end")
	return
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {
	size := len(bs)
	return size, nil
}
