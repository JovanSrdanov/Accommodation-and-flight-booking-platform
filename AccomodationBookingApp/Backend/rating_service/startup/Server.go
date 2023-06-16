package startup

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	rating "common/proto/rating_service/generated"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
	"log"
	"net"
	"rating_service/communication/handler"
	"rating_service/domain/service"
	"rating_service/persistence/repository"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server Server) Start() {
	neo4jClient := server.initNeo4jClient()
	ratingRepo := initRatingRepo(neo4jClient)
	ratingService := service.NewRatingService(*ratingRepo)
	ratingHandler := handler.NewRatingHandler(*ratingService)
	server.startGrpcServer(ratingHandler)
}

func (server Server) initNeo4jClient() neo4j.Driver {
	client, err := repository.GetClient(server.config.DbUri, server.config.DbUser, server.config.DbPass)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initRatingRepo(neo4jClient neo4j.Driver) *repository.RatingRepositoryNeo4J {
	repo, err := repository.NewRatingRepositoryNeo4J(neo4jClient)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
func (server Server) startGrpcServer(ratingHandler *handler.RatingHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	log.Println("port: " + server.config.Port)
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
	rating.RegisterRatingServiceServer(grpcServer, ratingHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/rating.RatingService/"

	return map[string][]model.Role{
		authServicePath + "RateAccommodation":            {model.Guest},
		authServicePath + "RateHost":                     {model.Guest},
		authServicePath + "DeleteRatingForAccommodation": {model.Guest},
		authServicePath + "DeleteRatingForHost":          {model.Guest},
	}
}
