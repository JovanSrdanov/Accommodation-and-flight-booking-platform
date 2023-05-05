package startup

import (
	"api_gateway/startup/configuration"
	authorization "common/proto/authorization_service/generated"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type Server struct {
	config  *configuration.Configuration
	handler *http.Handler
}

func NewServer(config *configuration.Configuration) *Server {
	mux := runtime.NewServeMux()
	handler := createAuthTokenMiddleware(mux)

	server := &Server{
		config:  config,
		handler: &handler,
	}
	server.initHandlers()
	// custom handlers with auth

	return server
}

func createAuthTokenMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		accessToken := authHeader[len("Bearer "):]
		ctx := context.WithValue(r.Context(), "access_token", accessToken)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	err := authorization.RegisterAuthorizationServiceHandlerFromEndpoint(context.TODO(), server.handler, authorizationEndpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), *server.handler))
}
