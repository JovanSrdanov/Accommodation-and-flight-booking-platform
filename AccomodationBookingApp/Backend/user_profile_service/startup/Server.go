package startup

import (
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
	grpcServer := grpc.NewServer()
	user_profile.RegisterUserProfileServiceServer(grpcServer, userProfileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
