package view

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	thirdservice "translasan-lite/ThirdService"
	"translasan-lite/common"
	"translasan-lite/core"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"google.golang.org/protobuf/proto"
)

type HttpServer struct {
	core.BaseComponent
	selfServer *gin.Engine
	// 未来可能会对selfServer的具体实现进行重新封装。嗯，大刀阔斧的事未来再说
	isDev bool
}

func (hs *HttpServer) ServiceID() int {
	return common.ServiceIdHttp
}

func (hs *HttpServer) Init(n *core.Node, cfg *core.ServiceConfig) {
	(&hs.BaseComponent).Init(n, cfg)
	hs.selfServer = gin.Default()

	// 创建一个CORS配置对象
	corsConfig := cors.DefaultConfig()
	// 如果你想允许任何源访问，你可以这样设置（不推荐用于生产环境）
	// corsConfig.AllowAllOrigins = true

	hs.isDev, _ = hs.Config.GetBool("is_dev")
	if hs.isDev {
		corsConfig.AllowOrigins = []string{"http://127.0.0.1:8080"}
		corsConfig.AllowMethods = []string{"GET", "POST"}
		corsConfig.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
		corsConfig.AllowCredentials = true
		// hs.selfServer.Use(Cors())
		// 使用CORS中间件
		hs.selfServer.Use(cors.New(corsConfig))
	}

	// 以上为跨域配置，实际上线记得去掉
	hs.selfServer.Use(hs.RecoveryMiddleware())
	hs.SetURLs()

}

func (hs *HttpServer) Run(ctx context.Context) error {

	addr, ok := hs.Config.GetString("http_addr")
	if !ok {
		// 默认地址如下
		addr = "127.0.0.1:8080"
	}

	hs.Log.Info("HttpServer starts running")
	go func() {
		hs.selfServer.Run(addr) // 监听并在 0.0.0.0:8080 上启动服务
	}()

	<-ctx.Done()
	hs.Log.Info("HttpServer stops running")

	return nil
}

func (hs *HttpServer) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录错误日志
				errorMsg := fmt.Sprintf("Recovered in %v: %v\n", c.Request.URL, r)
				hs.Log.Error(errorMsg)

				// 将 panic 转换为 Gin 可以处理的错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"message": fmt.Sprintf("%v", r),
				})

				// 阻止继续执行后续的中间件或处理程序
				return
			}
		}()

		// 继续执行后续的中间件或处理程序
		c.Next()
	}
}

func (hs *HttpServer) SetURLs() {
	type selfRouterConfig struct {
		path     string
		function gin.HandlerFunc
		method   string
	}
	SelfRouterConfigList := []*selfRouterConfig{
		{path: "/helloworld", function: HelloWorld, method: "GET"},
		{path: "/selectJpTable", function: hs.GetJpLiteTable, method: "POST"},
		{path: "/jplevelup", function: hs.RememberJpWord, method: "POST"},
		{path: "/jpleveldown", function: hs.ForgetJpWord, method: "POST"},
		{path: "/jpinsertupdate", function: hs.SaveJpWord, method: "POST"},
		{path: "/jpdelete", function: hs.DeleteJpWord, method: "POST"},
		{path: "/lookup", function: hs.TranslateJpWord, method: "POST"},
		{path: "/selectDetailTable", function: hs.GetJpDetailTable, method: "POST"},
		{path: "/getMarkdown", function: hs.GetMarkdownContent, method: "GET"},
		{path: "/saveMarkdown", function: hs.SaveMarkdownContent, method: "POST"},
		{path: "/askBotDemo", function: hs.AskBotDemo, method: "POST"},
	}
	var p string
	for _, config := range SelfRouterConfigList {
		if !hs.isDev {
			// 为了生产环境部署的改动
			p = "/back" + config.path
		} else {
			p = config.path
		}
		if config.method == "POST" {
			hs.selfServer.POST(p, config.function)
		} else if config.method == "GET" {
			hs.selfServer.GET(p, config.function)
		}
	}
}

func (hs *HttpServer) ReadProtoReq(r *http.Request, msg proto.Message) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return &HttpError{
			code: http.StatusBadRequest,
			err:  err,
		}
	}
	r.Body.Close()
	err = json.Unmarshal(data, msg)
	if err != nil {
		return &HttpError{
			code: http.StatusBadRequest,
			err:  err,
		}
	}
	return nil
}

func (hs *HttpServer) getThirdServiceClient() (thirdservice.IThirdService, error) {
	svc, ok := hs.FindService(common.ServiceIdThird)
	if !ok {
		return nil, common.ErrorInstance.ErrNoThirdService
	}
	ts, ok := svc.(thirdservice.IThirdService)
	if !ok {
		return nil, common.ErrorInstance.ErrInvalidThirdService
	}
	return ts, nil
}
