package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-leo/leo"
	"github.com/go-leo/leo/common/stringx"
	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	middlewarecontext "github.com/go-leo/leo/middleware/context"
	middlewarelog "github.com/go-leo/leo/middleware/log"
	"github.com/go-leo/leo/middleware/requestid"
	httpserver "github.com/go-leo/leo/runner/net/http/server"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())
	// 初始化app
	app := leo.NewApp(
		leo.Name("httpdemo"),
		// 日志打印
		leo.Logger(logger),
		leo.HTTP(&leo.HttpOptions{
			Port: 8080,
			// 全局中间件
			GinMiddlewares: []gin.HandlerFunc{
				requestid.GinMiddleware(),
				middlewarecontext.GinMiddleware(func(ctx context.Context) context.Context {
					traceID, _ := requestid.FromContext(ctx)
					return log.NewContext(ctx, logger.Clone().With(log.F{K: "TraceID", V: traceID}))
				}),
				middlewarelog.GinMiddleware(func(ctx context.Context) log.Logger { return log.FromContextOrDiscard(ctx) }),
			},
			Routers: []httpserver.Router{
				{
					HTTPMethod:   http.MethodPost,
					Path:         "/register",
					HandlerFuncs: []gin.HandlerFunc{Register},
				},
				{
					HTTPMethod:   http.MethodPost,
					Path:         "/login",
					HandlerFuncs: []gin.HandlerFunc{Login},
				},
			},
		}),
		leo.Management(&leo.ManagementOptions{
			Port: 16060,
		}),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

var users = map[string]*User{}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func Register(c *gin.Context) {
	u := new(User)
	err := c.BindJSON(u)
	if err != nil {
		_ = c.Error(err)
	}
	if stringx.IsBlank(u.Username) || stringx.IsBlank(u.Password) {
		c.String(http.StatusBadRequest, "username or password is blank")
		return
	}
	_, ok := users[u.Username]
	if ok {
		c.String(http.StatusBadRequest, "%s has exist")
		return
	}
	users[u.Username] = u
	c.String(http.StatusOK, "register success")
}

func Login(c *gin.Context) {
	u := new(User)
	err := c.BindJSON(u)
	if err != nil {
		_ = c.Error(err)
	}
	if stringx.IsBlank(u.Username) || stringx.IsBlank(u.Password) {
		c.String(http.StatusBadRequest, "username or password is blank")
		return
	}
	user, ok := users[u.Username]
	if !ok {
		c.String(http.StatusBadRequest, "%s not exist", u.Username)
		return
	}
	if user.Password != u.Password {
		c.String(http.StatusBadRequest, "password not right")
		return
	}
	c.String(http.StatusOK, "register success")
}
