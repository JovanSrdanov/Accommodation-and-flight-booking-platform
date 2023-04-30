package communication

import (
	"authorization_service/communication/handler"
	"authorization_service/configuration"
	"authorization_service/domain/service"
	"authorization_service/persistence/repository"
	authorization "common/proto/authorization_service/generated"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type Server struct {
	config *configuration.Configuration
}

func NewServer(config *configuration.Configuration) *Server {
	return &Server{
		config: config,
	}
}

func (server Server) Start() {
	postgresClient := server.initPostgresClient()
	accountCredentialsRepo := server.initAccountCredentialsRepo(postgresClient)
	accountCredentialsService := service.NewAccountCredentialsService(accountCredentialsRepo)
	accountCredentialsHandler := handler.NewAccountCredentialsHandler(accountCredentialsService)
	server.startGrpcServer(accountCredentialsHandler)
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

func (server Server) initAccountCredentialsRepo(postgresClient *gorm.DB) *repository.AccountCredentialsRepositoryPG {
	repo, err := repository.NewAccountCredentialsRepositoryPG(postgresClient)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func (server *Server) startGrpcServer(accountCredentialsHandler *handler.AccountCredentialsHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	authorization.RegisterAuthorizationServiceServer(grpcServer, accountCredentialsHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
