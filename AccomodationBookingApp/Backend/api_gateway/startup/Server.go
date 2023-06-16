package startup

import (
	"common/NotificationMessaging"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"context"
	"fmt"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"api_gateway/communication/handler"
	"api_gateway/communication/middleware"
	"authorization_service/domain/token"
	accommodation "common/proto/accommodation_service/generated"
	authorization "common/proto/authorization_service/generated"
	notification "common/proto/notification_service/generated"
	rating "common/proto/rating_service/generated"
	reservation "common/proto/reservation_service/generated"
	user_profile "common/proto/user_profile_service/generated"
)

type Server struct {
	config     *Configuration
	server     *gin.Engine
	subscriber messaging.Subscriber
}

const (
	QueueGroup = "api_gateway"
)

var wsConnections = make(map[uuid.UUID]*websocket.Conn) // Map to store websocket connections by user ID

func NewServer(config *Configuration) *Server {
	server := &Server{
		config: config,
		server: gin.Default(),
	}

	// OpenTelemetry
	server.server.Use(otelgin.Middleware("api-gateway"))

	prom := ginprometheus.NewPrometheus("gin")
	prom.Use(server.server)

	server.server.Use(middleware.AuthTokenParser())

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	})

	server.server.Use(corsMiddleware)

	server.server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Endpoint doesn't exist"})
	})

	grpcMux := runtime.NewServeMux()
	server.initGrpcHandlers(grpcMux)
	// Wrapping gin around runtime.ServeMux
	server.server.Group("/api-1/*any").Any("", gin.WrapH(grpcMux))

	server.initCustomHandlers(server.server.Group("/api-2"))

	server.server.GET("/ws", server.handleWebSocket)

	sendEventSubscriber := server.initSendEventSubscriber(server.config.SendNotificationToAPIGatewaySubject, QueueGroup)

	sendEventSubscriber.Subscribe(server.HandleMessages)

	return server
}

func (server *Server) HandleMessages(message *NotificationMessaging.NotificationMessage) {
	conn, ok := wsConnections[message.AccountID]
	if !ok {
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message.MessageForNotification))
	if err != nil {
		log.Println("Failed to send message through WebSocket connection:", err)
		return
	}
}

func (server *Server) initSendEventSubscriber(subject string, queueGroup string) messaging.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) handleWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	pasetoToken := c.Query("authorization")
	tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")
	tokenPayload, err := tokenMaker.VerifyToken(pasetoToken)
	userID := tokenPayload.ID
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	wsConnections[userID] = conn
	for {
		_, _, readErr := conn.ReadMessage()
		if readErr != nil {
			log.Println("WebSocket connection closed by client:", readErr)
			break
		}
	}
	delete(wsConnections, userID)
}

func (server *Server) initGrpcHandlers(mux *runtime.ServeMux) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.NewGRPUnaryClientInterceptor())}

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

	accommodationEndpoint := fmt.Sprintf("%s:%s", server.config.AccommodationHost, server.config.AccommodationPort)
	err = accommodation.RegisterAccommodationServiceHandlerFromEndpoint(context.TODO(), mux, accommodationEndpoint, opts)
	if err != nil {
		panic(err)
	}

	reservationEndpoint := fmt.Sprintf("%s:%s", server.config.ReservationHost, server.config.ReservationPort)
	err = reservation.RegisterReservationServiceHandlerFromEndpoint(context.TODO(), mux, reservationEndpoint, opts)
	if err != nil {
		panic(err)
	}

	notificationEndpoint := fmt.Sprintf("%s:%s", server.config.NotificationHost, server.config.NotificationPort)
	err = notification.RegisterNotificationServiceHandlerFromEndpoint(context.TODO(), mux, notificationEndpoint, opts)
	if err != nil {
		panic(err)
	}

	ratingEndpoint := fmt.Sprintf("%s:%s", server.config.RatingHost, server.config.RatingPort)
	err = rating.RegisterRatingServiceHandlerFromEndpoint(context.TODO(), mux, ratingEndpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (server *Server) initCustomHandlers(routerGroup *gin.RouterGroup) {
	tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")

	authorizationEndpoint := fmt.Sprintf("%s:%s", server.config.AuthorizationHost, server.config.AuthorizationPort)
	userProfileEndpoint := fmt.Sprintf("%s:%s", server.config.UserProfileHost, server.config.UserProfilePort)
	accommodationEndpoint := fmt.Sprintf("%s:%s", server.config.AccommodationHost, server.config.AccommodationPort)
	reservationEndpoint := fmt.Sprintf("%s:%s", server.config.ReservationHost, server.config.ReservationPort)
	notificationEndpoint := fmt.Sprintf("%s:%s", server.config.NotificationHost, server.config.NotificationPort)
	ratingEndpoint := fmt.Sprintf("%s:%s", server.config.RatingHost, server.config.RatingPort)

	userInfoHandler := handler.NewUserHandler(authorizationEndpoint, userProfileEndpoint, notificationEndpoint, tokenMaker)
	userInfoHandler.Init(routerGroup)

	accommodationHandler := handler.NewAccommodationHandler(accommodationEndpoint, reservationEndpoint, authorizationEndpoint, userProfileEndpoint, ratingEndpoint, tokenMaker)
	accommodationHandler.Init(routerGroup)
}

func (server *Server) Start() {
	port := fmt.Sprintf(":%s", server.config.Port)
	err := server.server.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
