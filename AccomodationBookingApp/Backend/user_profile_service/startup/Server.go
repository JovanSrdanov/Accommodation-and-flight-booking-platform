package startup

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	user_profile "common/proto/user_profile_service/generated"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"user_profile_service/communication/handler"
	"user_profile_service/domain/service"
	"user_profile_service/persistence/repository"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server Server) Start() {
	postgresClient := server.initPostgresClient()
	userProfileRepo := initUserProfileRepo(postgresClient)
	userProfileService := service.NewUserProfileService(*userProfileRepo)
	userProfileHandler := handler.NewUserProfileHandler(*userProfileService)
	server.startGrpcServer(userProfileHandler)
}

func initUserProfileRepo(client *gorm.DB) *repository.UserProfileRepositoryPG {
	repo, err := repository.NewUserProfileRepositoryPG(client)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
func (server Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetClient(
		server.config.DBHost, server.config.DBUser,
		server.config.DBPass, server.config.DBName,
		server.config.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
func (server Server) startGrpcServer(userProfileHandler *handler.UserProfileHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")
	protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	authInterceptor := interceptor.NewAuthServerInterceptor(tokenMaker, protectedMethodsWithAllowedRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	user_profile.RegisterUserProfileServiceServer(grpcServer, userProfileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/user_profile.UserProfileService/"

	return map[string][]model.Role{
		authServicePath + "Update": {model.Guest},
	}
}
