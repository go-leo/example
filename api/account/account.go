package account

import "context"

//go:generate gors -service AccountServer

// AccountServer
// @GORS @Path(/account)
type AccountServer interface {
	// Register
	// @GORS @POST @Path(/register) @JSONBinding @JSONRender
	Register(ctx context.Context, user *User) (*Empty, error)
	// Login
	// @GORS @POST @Path(/login) @JSONBinding @JSONRender
	Login(ctx context.Context, user *User) (*Empty, error)
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Empty struct{}
