package startup

import (
	"api_gateway/communication/handler"
	"api_gateway/communication/middleware"
	"authorization_service/domain/token"
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Server struct {
	config *Configuration
	server *gin.Engine
}

func NewServer(config *Configuration) *Server {
	server := &Server{
		config: config,
		server: gin.Default(),
	}

	server.server.Use(middleware.AuthTokenParser())

	server.server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Endpoint doesn't exist"})
	})

	grpcMux := runtime.NewServeMux()
	server.initGrpcHandlers(grpcMux)
	//Wrapping gin around runtime.ServeMux
	server.server.Group("/api-1/*any").Any("", gin.WrapH(grpcMux))

	server.initCustomHandlers(server.server.Group("/api-2"))

	return server
}

func (server *Server) initGrpcHandlers(mux *runtime.ServeMux) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	err := authorization.RegisterAuthorizationServiceHandlerFromEndpoint(context.TODO(), mux, authorizationEndpoint, opts)
	if err != nil {
		panic(err)
	}

	userProfileEndpoint := fmt.Sprintf("%s:%s", server.config.UserProfileHost, server.config.UserProfilePort)
	err = user_profile.RegisterUserProfileServiceHandlerFromEndpoint(context.TODO(), mux, userProfileEndpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (server *Server) initCustomHandlers(routerGroup *gin.RouterGroup) {
	tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")

	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	userProfileEndpoint := fmt.Sprintf("%s:%s", server.config.UserProfileHost, server.config.UserProfilePort)

	userInfoHandler := handler.NewUserHandler(authorizationEndpoint, userProfileEndpoint, tokenMaker)
	userInfoHandler.Init(routerGroup)
}

func (server *Server) Start() {
	port := fmt.Sprintf(":%s", server.config.Port)
	err := server.server.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
