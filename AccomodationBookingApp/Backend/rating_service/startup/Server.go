package startup

import (
	rating "common/proto/rating_service/generated"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"rating_service/communication/handler"
	"rating_service/communication/middleware"
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
	sendEventPublisher := server.initSendEventPublisher(server.config.SendEventToNotificationServiceSubject)
	neo4jClient := server.initNeo4jClient()
	ratingRepo := initRatingRepo(neo4jClient, sendEventPublisher)
	ratingService := service.NewRatingService(*ratingRepo, sendEventPublisher)
	ratingHandler := handler.NewRatingHandler(*ratingService, sendEventPublisher)
	server.startGrpcServer(ratingHandler)
}

func (server *Server) initSendEventPublisher(subject string) messaging.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server Server) initNeo4jClient() neo4j.Driver {
	client, err := repository.GetClient(server.config.DbUri, server.config.DbUser, server.config.DbPass)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initRatingRepo(neo4jClient neo4j.Driver, publisher messaging.Publisher) *repository.RatingRepositoryNeo4J {
	repo, err := repository.NewRatingRepositoryNeo4J(neo4jClient, publisher)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/root/cert/rating-service-cert.pem", "/root/cert/rating-service-key.pem")
	if err != nil {
		return nil, err
	}

	// Load certificate of the CA who signed the certificate
	pemServerCA, err := ioutil.ReadFile("/root/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add the CA certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func (server Server) startGrpcServer(ratingHandler *handler.RatingHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	log.Println("port: " + server.config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")
	//protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	//authInterceptor := interceptor.NewAuthServerInterceptor(tokenMaker, protectedMethodsWithAllowedRoles)

	// Enable TLS
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(middleware.NewGRPUnaryServerInterceptor()), /*authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()*/
	)
	rating.RegisterRatingServiceServer(grpcServer, ratingHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//// returns a map which consists of a list of grpc methods and allowed roles for each of them
//func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
//	const authServicePath = "/rating.RatingService/"
//
//	return map[string][]model.Role{
//		authServicePath + "RateAccommodation":            {model.Guest},
//		authServicePath + "RateHost":                     {model.Guest},
//		authServicePath + "DeleteRatingForAccommodation": {model.Guest},
//		authServicePath + "DeleteRatingForHost":          {model.Guest},
//	}
//}
