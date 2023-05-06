package startup

import (
	"api_gateway/communication/handler"
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
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

	server.server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Endpoint doesn't exist"})
	})

	grpcMux := runtime.NewServeMux()
	server.initGrpcHandlers(grpcMux)
	//Wrapping gin around runtime.ServeMux
	server.server.Group("/api-1/*any").Any("", gin.WrapH(grpcMux))

	server.initCustomHandlers(server.server.Group("/api-2"))

	//When it initializes all handlers on basic mux, we wrap it in a middleware(handler)

	// custom handlers with auth
	//TODO better name
	//TODO GIN MIDLEWARE
	//authHandler := createAuthTokenMiddleware(server.mux)
	//server.handler = &authHandler

	return server
}

func createAuthTokenMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		ctx := request.Context()
		// if authorization header is provided, embeds the token to the context which is being sent to a grpc server,
		// else just sends the default context
		if authHeader != "" {
			accessToken := authHeader[len("Bearer "):]
			ctx := context.WithValue(ctx, "Authorization", accessToken)
			//ctx := metadata.AppendToOutgoingContext(ctx, "Authorization", accessToken)
			handler.ServeHTTP(writer, request.WithContext(ctx))
			return
		}
		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
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
	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	userProfileEndpoint := fmt.Sprintf("%s:%s", server.config.UserProfileHost, server.config.UserProfilePort)

	userInfoHandler := handler.NewUserHandler(authorizationEndpoint, userProfileEndpoint)
	userInfoHandler.Init(routerGroup)
}

func (server *Server) Start() {
	port := fmt.Sprintf(":%s", server.config.Port)
	err := server.server.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
