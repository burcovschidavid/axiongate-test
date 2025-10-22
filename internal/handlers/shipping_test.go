package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"shipping-api/internal/core/domain"
	"shipping-api/internal/core/service"
	"shipping-api/internal/testutil"
	"strings"
	"testing"
)

func TestShippingHandler_CreateShipment_SingleProvider(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)

	mockProvider := testutil.NewMockShippingProvider("A", "http://a.local")
	shippingService.RegisterProvider(mockProvider)

	handler := NewShippingHandler(shippingService)

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping?provider=A", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
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
}

func TestShippingHandler_CreateShipment_Broadcast(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)

	providerA := testutil.NewMockShippingProvider("A", "http://a.local")
	providerB := testutil.NewMockShippingProvider("B", "http://b.local")
	shippingService.RegisterProvider(providerA)
	shippingService.RegisterProvider(providerB)

	handler := NewShippingHandler(shippingService)

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var responses []*domain.ShipmentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &responses); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}

	providerNames := make(map[string]bool)
	for _, resp := range responses {
		providerNames[resp.Provider] = true
	}

	if !providerNames["A"] || !providerNames["B"] {
		t.Error("expected responses from both providers A and B")
	}
}

func TestShippingHandler_CreateShipment_InvalidJSON(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)
	handler := NewShippingHandler(shippingService)

	invalidJSON := []byte(`{"invalid": json}`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping?provider=A", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", w.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errorResponse["error"] == "" {
		t.Error("expected error message in response")
	}
}

func TestShippingHandler_CreateShipment_MethodNotAllowed(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)
	handler := NewShippingHandler(shippingService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/createShipping", nil)
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status code 405, got %d", w.Code)
	}
}

func TestShippingHandler_CreateShipment_ProviderNotFound(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	shippingService := service.NewShippingService(mockRepo)
	handler := NewShippingHandler(shippingService)

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping?provider=NonExistent", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if !strings.Contains(errorResponse["error"], "not found") {
		t.Errorf("expected 'not found' in error message, got %s", errorResponse["error"])
	}
}

type mockFailingService struct{}

func (m *mockFailingService) ProcessShipment(ctx context.Context, request *domain.GenericShippingRequest, providerName string) (*domain.ShipmentResponse, error) {
	return nil, errors.New("service error")
}

func (m *mockFailingService) BroadcastShipment(ctx context.Context, request *domain.GenericShippingRequest) ([]*domain.ShipmentResponse, error) {
	return nil, errors.New("broadcast error")
}

func TestShippingHandler_CreateShipment_ServiceError(t *testing.T) {
	handler := NewShippingHandler(&mockFailingService{})

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping?provider=A", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}
}

func TestShippingHandler_CreateShipment_BroadcastError(t *testing.T) {
	handler := NewShippingHandler(&mockFailingService{})

	request := testutil.CreateSampleShippingRequest()
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/createShipping", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateShipment(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}
}
