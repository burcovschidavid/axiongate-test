package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type MockResponse struct {
	TrackingID string                 `json:"trackingId"`
	AWB        string                 `json:"awb"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Provider   string                 `json:"provider"`
	Timestamp  string                 `json:"timestamp"`
	Request    map[string]interface{} `json:"request,omitempty"`
}

func main() {
	provider := os.Getenv("PROVIDER_NAME")
	if provider == "" {
		provider = "UNKNOWN"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/createShipping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request", http.StatusBadRequest)
			return
		}

		var requestData map[string]interface{}
		json.Unmarshal(body, &requestData)

		trackingPrefix := fmt.Sprintf("%s-TRACK", provider)
		awbPrefix := fmt.Sprintf("%s-AWB", provider)

		response := MockResponse{
			TrackingID: fmt.Sprintf("%s-%d", trackingPrefix, rand.Intn(999999)),
			AWB:        fmt.Sprintf("%s-%d", awbPrefix, rand.Intn(999999)),
			Status:     "success",
			Message:    fmt.Sprintf("Shipment created successfully via provider %s", provider),
			Provider:   provider,
			Timestamp:  time.Now().Format(time.RFC3339),
			Request:    requestData,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		log.Printf("[%s] Received shipment request - Tracking: %s", provider, response.TrackingID)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Mock provider %s server starting on port %s", provider, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
