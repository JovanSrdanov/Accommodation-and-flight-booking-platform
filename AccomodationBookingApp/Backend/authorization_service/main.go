package main

import (
	"authorization_service/communication/middleware"
	"authorization_service/startup"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"os"
)

func main() {
	// OpenTelemetry
	tp, err := middleware.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	configuration := startup.NewConfiguration()
	server := startup.NewServer(configuration)
	log.Printf("Authorization service started, running on %s:%s", os.Getenv("AUTHORIZATION_SERVICE_HOST"), os.Getenv("AUTHORIZATION_SERVICE_PORT"))
	server.Start()
}
