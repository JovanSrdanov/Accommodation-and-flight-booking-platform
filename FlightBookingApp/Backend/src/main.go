package main

import (
	"FlightBookingApp/docs"
	"FlightBookingApp/endpoints"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO Aleksandar: namestiti kako treba, kada se uvede autorizacija i dodati tagove za autorizaciju na svaki endpoint
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
func main() {
	logger := log.New(os.Stdout, "[flight-app-api] ", log.LstdFlags)

	//DB initialization
	dbLogger := log.New(os.Stdout, "[mongo-db] ", log.LstdFlags)

	dbClient, err := GetClient(dbLogger)
	if err != nil {
		logger.Println(err.Error())
	}

	Connect(context.Background(), dbClient, dbLogger)
	defer Disconnect(context.Background(), dbClient, dbLogger)

	//Setting swagger runtime variables
	port := ":" + os.Getenv("PORT")
	//port := ":4200"
	docs.SwaggerInfo.Title = "Flights booking app"
	docs.SwaggerInfo.Host = "localhost" + port
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	//Routes definition
	//Remove debug logs
	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	apiRoutes := router.Group(docs.SwaggerInfo.BasePath)
	{
		//endpoints.DefineFlightEndpoints(apiRoutes, dbClient)
		endpoints.DefineAccountEndpoints(apiRoutes, dbClient)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Server initialization
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		// service connections
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error: %s\n", err)
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//Here it waits for interrupt
	logger.Printf("Listening on port %s\n", port)
	<-quit

	log.Println("Shutdown Server...")

	// The context is used to inform the server it has n seconds to finish
	// the request it is currently handling
	// if n seconds pass cancel function will be called
	timeout, _ := strconv.ParseInt(os.Getenv("GRACEFUL_SHUTDOWN_TIMEOUT"), 10, 32)
	//timeout := 2

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of n seconds.
	// Waits for context to say that n second timeout passed
	select {
	case <-ctx.Done():
		log.Printf("Timeout of %d seconds passed.", timeout)
	}

	logger.Println("Server exiting")
}
