package client

import (
	"accommodation_service/communication/middleware"
	authorization "common/proto/authorization_service/generated"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed the certificates
	pemServerCA, err := ioutil.ReadFile("/root/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("/root/cert/accommodation-service-cert.pem", "/root/cert/accommodation-service-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func NewAccountCredentialsClient(address string) authorization.AuthorizationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to AccountCredentials service: %v", err)
	}

	return authorization.NewAuthorizationServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	return grpc.Dial(address, grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(middleware.NewGRPUnaryClientInterceptor()))
}
