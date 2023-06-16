package main

import (
	"api_gateway/communication/middleware"
	"api_gateway/startup"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
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

	config := startup.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
