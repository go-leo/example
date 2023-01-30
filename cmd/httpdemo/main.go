package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-leo/gors"
	"github.com/go-leo/leo/v2"
	leohttp "github.com/go-leo/leo/v2/http"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	"github.com/go-leo/stringx"

	"github.com/go-leo/example/v2/api/account"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())
	// 初始化app
	engine := gors.AppendRoutes(gin.New(), account.AccountServerRoutes(new(AccountService))...)
	httpSrv, err := leohttp.NewServer(8080, engine)
	if err != nil {
		panic(err)
	}
	app := leo.NewApp(
		leo.Name("httpdemo"),
		// 日志打印
		leo.Logger(logger),
		leo.HTTP(httpSrv),
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

type AccountService struct{}

func (a *AccountService) Register(ctx context.Context, u *account.User) (*account.Empty, error) {
	if stringx.IsBlank(u.Username) || stringx.IsBlank(u.Password) {
		return nil, errors.New("username or password is blank")
	}
	_, ok := users[u.Username]
	if ok {
		return nil, fmt.Errorf("%s exist", u.Username)
	}
	users[u.Username] = &User{
		Username: u.Username,
		Password: u.Password,
	}
	return &account.Empty{}, nil
}

func (a *AccountService) Login(ctx context.Context, u *account.User) (*account.Empty, error) {
	if stringx.IsBlank(u.Username) || stringx.IsBlank(u.Password) {
		return nil, errors.New("username or password is blank")
	}
	user, ok := users[u.Username]
	if !ok {
		return nil, fmt.Errorf("%s not exist", u.Username)
	}
	if user.Password != u.Password {
		return nil, errors.New("password not right")
	}
	return &account.Empty{}, nil
}
