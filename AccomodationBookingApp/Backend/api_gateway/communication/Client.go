package communication

import (
	accommodation "common/proto/accommodation_service/generated"
	authorization "common/proto/authorization_service/generated"
	notification "common/proto/notification_service/generated"
	reservation "common/proto/reservation_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func NewAuthorizationClient(address string) authorization.AuthorizationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Authorization service: %v", err)
	}
	return authorization.NewAuthorizationServiceClient(conn)
}

func NewUserProfileClient(address string) user_profile.UserProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to UserProfile service: %v", err)
	}
	return user_profile.NewUserProfileServiceClient(conn)
}

func NewAccommodationClient(address string) accommodation.AccommodationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to UserProfile service: %v", err)
	}
	return accommodation.NewAccommodationServiceClient(conn)
}

func NewReservationClient(address string) reservation.ReservationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to UserProfile service: %v", err)
	}
	return reservation.NewReservationServiceClient(conn)
}

func NewNotificationClient(address string) notification.NotificationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to UserProfile service: %v", err)
	}
	return notification.NewNotificationServiceClient(conn)
}
