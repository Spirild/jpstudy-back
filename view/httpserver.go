package view

import (
	"context"
	"translasan-lite/core"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	core.BaseComponent
	selfServer *gin.Engine
	// 未来可能会对selfServer的具体实现进行重新封装。嗯，大刀阔斧的事未来再说
}

func (hs *HttpServer) Init(n *core.Node, cfg *core.ServiceConfig) {
	(&hs.BaseComponent).Init(n, cfg)
	hs.selfServer = gin.Default()
	hs.SetURLs()
}

func (hs *HttpServer) Run(ctx context.Context) error {

	addr, ok := hs.Config.GetString("http_addr")
	if !ok {
		// 默认地址如下
		addr = "127.0.0.1:8080"
	}
	hs.selfServer.Run(addr) // 监听并在 0.0.0.0:8080 上启动服务
	<-ctx.Done()
	hs.Log.Info("This Server starts running")

	return nil
}

func (hw *HttpServer) SetURLs() {
	hw.selfServer.GET("/helloworld", HelloWorld)
	// 注册函数可能还要专门写一个地方
}
