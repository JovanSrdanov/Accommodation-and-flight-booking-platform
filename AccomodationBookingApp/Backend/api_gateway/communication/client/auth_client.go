package client

import (
	authorization "common/proto/authorization_service/generated"
	"context"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	service  authorization.AuthorizationServiceClient
	username string
	password string
}

func NewAuthClient(conn *grpc.ClientConn, username string, password string) *AuthClient {
	service := authorization.NewAuthorizationServiceClient(conn)
	return &AuthClient{service, username, password}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &authorization.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
