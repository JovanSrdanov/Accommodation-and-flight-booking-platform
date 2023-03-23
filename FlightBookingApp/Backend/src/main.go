package main

import (
	"FlightBookingApp/endpoints"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "[flight-app-api] ", log.LstdFlags)

	//Db initialization (cant be extracted because of defer Disconnect has to be in main
	dbClient, err := GetClient()
	if err != nil {
		logger.Println(err.Error())
	}
	dbLogger := log.New(os.Stdout, "[mongo-db] ", log.LstdFlags)

	Connect(context.Background(), dbClient, dbLogger)
	defer Disconnect(context.Background(), dbClient, dbLogger)

	//Routes definition
	router := gin.Default()
	apiRoutes := router.Group("/api")
	{
		endpoints.DefineFlightEndpoints(apiRoutes, dbClient)
	}

	//Server initialization

	port := ":" + os.Getenv("PORT")
	//port := ":4200"

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
	logger.Printf("Listening on port %s\n", port)

	//Graceful shutdown
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//Here it waits for interrupt
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
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 2 seconds.")
	}

	logger.Println("Server exiting")
}
