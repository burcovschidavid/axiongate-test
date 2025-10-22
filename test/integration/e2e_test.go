package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shipping-api/internal/adapters/providers/providerA"
	"shipping-api/internal/adapters/providers/providerB"
	"shipping-api/internal/core/domain"
	"shipping-api/internal/core/service"
	"shipping-api/internal/handlers"
	"shipping-api/internal/testutil"
	"testing"
)

func setupMockProviderServers() (*httptest.Server, *httptest.Server) {
	providerAServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"trackingId": "A-TRACK-123",
			"awb":        "A-AWB-123",
			"message":    "Shipment created successfully",
			"status":     "success",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))

	providerBServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"trackingId": "B-TRACK-456",
			"awb":        "B-AWB-456",
			"message":    "Shipment accepted",
			"status":     "success",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))

	return providerAServer, providerBServer
}

func TestE2E_CreateShipment_ProviderA(t *testing.T) {
	providerAServer, providerBServer := setupMockProviderServers()
	defer providerAServer.Close()
	defer providerBServer.Close()

	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)

	providerAAdapter := providerA.NewAdapter(providerAServer.URL)
	providerBAdapter := providerB.NewAdapter(providerBServer.URL)
	shippingService.RegisterProvider(providerAAdapter)
	shippingService.RegisterProvider(providerBAdapter)

	handler := handlers.NewShippingHandler(shippingService)

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping?provider=A", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response domain.ShipmentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Provider != "A" {
		t.Errorf("expected provider A, got %s", response.Provider)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.TrackingID != "A-TRACK-123" {
		t.Errorf("expected tracking ID A-TRACK-123, got %s", response.TrackingID)
	}

	if mockRepo.GetRecordCount() != 1 {
		t.Errorf("expected 1 record in database, got %d", mockRepo.GetRecordCount())
	}
}

func TestE2E_CreateShipment_Broadcast(t *testing.T) {
	providerAServer, providerBServer := setupMockProviderServers()
	defer providerAServer.Close()
	defer providerBServer.Close()

	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)

	providerAAdapter := providerA.NewAdapter(providerAServer.URL)
	providerBAdapter := providerB.NewAdapter(providerBServer.URL)
	shippingService.RegisterProvider(providerAAdapter)
	shippingService.RegisterProvider(providerBAdapter)

	handler := handlers.NewShippingHandler(shippingService)

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var responses []*domain.ShipmentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &responses); err != nil {
		t.Fatalf("failed to unmarshal responses: %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}

	providerMap := make(map[string]*domain.ShipmentResponse)
	for _, resp := range responses {
		providerMap[resp.Provider] = resp
		if !resp.Success {
			t.Errorf("expected success for provider %s", resp.Provider)
		}
	}

	if providerMap["A"] == nil {
		t.Error("missing response from provider A")
	}

	if providerMap["B"] == nil {
		t.Error("missing response from provider B")
	}

	if mockRepo.GetRecordCount() != 2 {
		t.Errorf("expected 2 records in database, got %d", mockRepo.GetRecordCount())
	}
}

func TestE2E_CreateShipment_ProviderFailure(t *testing.T) {
	providerAServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer providerAServer.Close()

	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)

	providerAAdapter := providerA.NewAdapter(providerAServer.URL)
	shippingService.RegisterProvider(providerAAdapter)

	handler := handlers.NewShippingHandler(shippingService)

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping?provider=A", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response domain.ShipmentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Success {
		t.Error("expected success to be false for failed provider")
	}

	if mockRepo.GetRecordCount() != 0 {
		t.Errorf("expected 0 records for failed shipment, got %d", mockRepo.GetRecordCount())
	}
}

func TestE2E_MapperIntegration_ProviderA(t *testing.T) {
	providerAServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var receivedRequest providerA.Request
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("provider A received invalid JSON: %v", err)
		}

		if receivedRequest.Weight.Value != 1000 {
			t.Errorf("expected weight 1000, got %f", receivedRequest.Weight.Value)
		}

		if receivedRequest.Shipper.Contact.Name != "ABC Associates" {
			t.Errorf("expected shipper name ABC Associates, got %s", receivedRequest.Shipper.Contact.Name)
		}

		if receivedRequest.PrintType != "AWBOnly" {
			t.Errorf("expected print type AWBOnly, got %s", receivedRequest.PrintType)
		}

		if len(receivedRequest.CustomsDeclarations) != 1 {
			t.Errorf("expected 1 customs declaration, got %d", len(receivedRequest.CustomsDeclarations))
		}

		response := map[string]interface{}{
			"trackingId": "TEST-123",
			"status":     "success",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer providerAServer.Close()

	adapter := providerA.NewAdapter(providerAServer.URL)
	request := testutil.CreateSampleShippingRequest()

	response, err := adapter.CreateShipment(context.Background(), request)

	if err != nil {
		t.Fatalf("failed to create shipment: %v", err)
	}

	if !response.Success {
		t.Error("expected successful response")
	}
}

func TestE2E_MapperIntegration_ProviderB(t *testing.T) {
	providerBServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var receivedRequest providerB.Request
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("provider B received invalid JSON: %v", err)
		}

		if receivedRequest.Weight != 1.0 {
			t.Errorf("expected weight 1.0 kg, got %f", receivedRequest.Weight)
		}

		if receivedRequest.Shipper != "ABC Associates" {
			t.Errorf("expected shipper ABC Associates, got %s", receivedRequest.Shipper)
		}

		if receivedRequest.ShipperCPerson != "ABC Associates" {
			t.Errorf("expected shipper contact ABC Associates, got %s", receivedRequest.ShipperCPerson)
		}

		if len(receivedRequest.PackageRequest) != 1 {
			t.Errorf("expected 1 package, got %d", len(receivedRequest.PackageRequest))
		}

		if len(receivedRequest.ExportItemDeclarationRequest) != 1 {
			t.Errorf("expected 1 export item, got %d", len(receivedRequest.ExportItemDeclarationRequest))
		}

		response := map[string]interface{}{
			"trackingId": "TEST-456",
			"status":     "success",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer providerBServer.Close()

	adapter := providerB.NewAdapter(providerBServer.URL)
	request := testutil.CreateSampleShippingRequest()

	response, err := adapter.CreateShipment(context.Background(), request)

	if err != nil {
		t.Fatalf("failed to create shipment: %v", err)
	}

	if !response.Success {
		t.Error("expected successful response")
	}
}

func TestE2E_ConcurrentRequests(t *testing.T) {
	providerAServer, providerBServer := setupMockProviderServers()
	defer providerAServer.Close()
	defer providerBServer.Close()

	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)

	providerAAdapter := providerA.NewAdapter(providerAServer.URL)
	providerBAdapter := providerB.NewAdapter(providerBServer.URL)
	shippingService.RegisterProvider(providerAAdapter)
	shippingService.RegisterProvider(providerBAdapter)

	handler := handlers.NewShippingHandler(shippingService)

	concurrentRequests := 10

	results := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			request := testutil.CreateSampleShippingRequest()
			requestBody, _ := json.Marshal(request)

			provider := "A"
			if index%2 == 0 {
				provider = "B"
			}

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/api/v1/createShipping?provider=%s", provider),
				bytes.NewBuffer(requestBody),
			)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateShipment(w, req)

			if w.Code != http.StatusOK {
				results <- fmt.Errorf("request %d failed with status %d", index, w.Code)
				return
			}

			results <- nil
		}(i)
	}

	for i := 0; i < concurrentRequests; i++ {
		if err := <-results; err != nil {
			t.Error(err)
		}
	}

	if mockRepo.GetRecordCount() != concurrentRequests {
		t.Errorf("expected %d records, got %d", concurrentRequests, mockRepo.GetRecordCount())
	}
}
