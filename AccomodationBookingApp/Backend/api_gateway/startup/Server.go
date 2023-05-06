package startup

import (
	"api_gateway/communication/handler"
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type Server struct {
	config  *Configuration
	mux     *runtime.ServeMux
	handler *http.Handler
}

func NewServer(config *Configuration) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()

	//When it initializes all handlers on basic mux, we wrap it in middleware(handler)

	// custom handlers with auth
	//TODO better name
	authHandler := createAuthTokenMiddleware(server.mux)
	server.handler = &authHandler
	return server
}

func createAuthTokenMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		ctx := request.Context()
		if authHeader != "" {
			accessToken := authHeader[len("Bearer "):]
			ctx := context.WithValue(ctx, "access_token", accessToken)
			handler.ServeHTTP(writer, request.WithContext(ctx))
			return
		}
		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	err := authorization.RegisterAuthorizationServiceHandlerFromEndpoint(context.TODO(), server.mux, authorizationEndpoint, opts)
	if err != nil {
		panic(err)
	}

	userProfileEndpoint := fmt.Sprintf("%s:%s", server.config.UserProfileHost, server.config.UserProfilePort)
	err = user_profile.RegisterUserProfileServiceHandlerFromEndpoint(context.TODO(), server.mux, userProfileEndpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (server *Server) initCustomHandlers() {
	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	userProfileEndpoint := fmt.Sprintf("%s:%s", server.config.UserProfileHost, server.config.UserProfilePort)

	userInfoHandler := handler.NewUserInfoHandler(authorizationEndpoint, userProfileEndpoint)
	userInfoHandler.Init(server.mux)
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), *server.handler))
}
