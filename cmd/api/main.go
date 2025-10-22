package main

import (
	"fmt"
	"log"
	"net/http"
	"shipping-api/internal/adapters/providers/providerA"
	"shipping-api/internal/adapters/providers/providerB"
	"shipping-api/internal/adapters/repository"
	"shipping-api/internal/core/service"
	"shipping-api/internal/handlers"
	"shipping-api/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()

	shippingService := service.NewShippingService(repo)

	providerAAdapter := providerA.NewAdapter(cfg.ProviderAURL)
	shippingService.RegisterProvider(providerAAdapter)

	providerBAdapter := providerB.NewAdapter(cfg.ProviderBURL)
	shippingService.RegisterProvider(providerBAdapter)

	handler := handlers.NewShippingHandler(shippingService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/createShipping", handler.CreateShipment)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("server starting on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
