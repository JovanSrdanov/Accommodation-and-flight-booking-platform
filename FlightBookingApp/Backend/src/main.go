package main

import (
	"FlightBookingApp/dependencyInjection"
	"FlightBookingApp/docs"
	"FlightBookingApp/endpoints"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO : namestiti kako treba, kada se uvede autorizacija i dodati tagove za autorizaciju na svaki endpoint
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOriginFunc = func(origin string) bool { return true }
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))
	depContainer := dependencyInjection.NewDependencyContainer()

	apiRoutes := router.Group(docs.SwaggerInfo.BasePath)
	{
		endpoints.DefineAccountEndpoints(apiRoutes, dbClient, depContainer)
		endpoints.DefineUserEndpoints(apiRoutes, dbClient, depContainer)

		RegisterJwtBundle(depContainer)

		endpoints.DefineFlightEndpoints(apiRoutes, dbClient, depContainer)
		endpoints.DefineAirportEndpoints(apiRoutes, dbClient, depContainer)
		endpoints.DefineTicketEndpoints(apiRoutes, dbClient, depContainer)
	}
	depContainer.PrintAllDependencies()
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

	DrawGoGopher()

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

func DrawGoGopher() {
	cyan := color.New(color.FgCyan).PrintlnFunc()
	cyan("------------------------------------------------------------------------------------")
	cyan("         ,_---~~~~~----._         ")
	cyan("  _,,_,*^____      _____``*g*\\\"*, ")
	cyan(" / __/ /'     ^.  /      \\ ^@q   f ")
	cyan("[  @f | @))    |  | @))   l  0 _/  ")
	cyan(" \\`/   \\~____ / __ \\_____/    \\   ")
	cyan("  |           _l__l_           I   ")
	cyan("  }          [______]           I  \t GOLANG! GO! GO! GO!")
	cyan("  ]            | | |            |  ")
	cyan("  ]             ~ ~             |  \t YOU CAN DO IT!")
	cyan("  |                            |   \t THE GO GOPHER BELIEVES IN YOU!")
	cyan("   |                           |   \t DON'T STOP WHEN YOU ARE TIRED. STOP WHEN YOU ARE DONE.")
	cyan("   |                           |   ")
	cyan("   |                           |   ")
	cyan("------------------------------------------------------------------------------------")
}

func RegisterJwtBundle(container *dependencyInjection.DependencyContainer) {
	accountRepo := container.GetRepository("account").(repository.AccountRepository)
	userRepo := container.GetRepository("user").(repository.UserRepository)
	jwtServ := service.NewJwtService(accountRepo, userRepo)

	container.RegisterService("jwt", jwtServ)
}
