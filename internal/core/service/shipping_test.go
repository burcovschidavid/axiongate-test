package service

import (
	"context"
	"errors"
	"shipping-api/internal/core/domain"
	"shipping-api/internal/testutil"
	"testing"
)

func TestShippingService_ProcessShipment_Success(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	service := NewShippingService(mockRepo)

	mockProvider := testutil.NewMockShippingProvider("TestProvider", "http://test.local")
	service.RegisterProvider(mockProvider)

	request := testutil.CreateSampleShippingRequest()

	response, err := service.ProcessShipment(context.Background(), request, "TestProvider")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.Provider != "TestProvider" {
		t.Errorf("expected provider TestProvider, got %s", response.Provider)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if mockRepo.GetRecordCount() != 1 {
		t.Errorf("expected 1 record in repository, got %d", mockRepo.GetRecordCount())
	}
}

func TestShippingService_ProcessShipment_ProviderNotFound(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	service := NewShippingService(mockRepo)

	request := testutil.CreateSampleShippingRequest()

	_, err := service.ProcessShipment(context.Background(), request, "NonExistentProvider")

	if err == nil {
		t.Fatal("expected error for non-existent provider, got nil")
	}

	if mockRepo.GetRecordCount() != 0 {
		t.Errorf("expected 0 records in repository, got %d", mockRepo.GetRecordCount())
	}
}

func TestShippingService_ProcessShipment_ProviderError(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	service := NewShippingService(mockRepo)

	mockProvider := testutil.NewMockShippingProvider("ErrorProvider", "http://test.local")
	mockProvider.SetCreateShipmentFunc(func(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error) {
		return nil, errors.New("provider error")
	})
	service.RegisterProvider(mockProvider)

	request := testutil.CreateSampleShippingRequest()

	_, err := service.ProcessShipment(context.Background(), request, "ErrorProvider")

	if err == nil {
		t.Fatal("expected error from provider, got nil")
	}

	if mockRepo.GetRecordCount() != 0 {
		t.Errorf("expected 0 records in repository for failed request, got %d", mockRepo.GetRecordCount())
	}
}

func TestShippingService_BroadcastShipment_Success(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	service := NewShippingService(mockRepo)

	providerA := testutil.NewMockShippingProvider("A", "http://a.local")
	providerB := testutil.NewMockShippingProvider("B", "http://b.local")
	service.RegisterProvider(providerA)
	service.RegisterProvider(providerB)

	request := testutil.CreateSampleShippingRequest()

	responses, err := service.BroadcastShipment(context.Background(), request)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}

	providerNames := make(map[string]bool)
	for _, resp := range responses {
		providerNames[resp.Provider] = true
		if !resp.Success {
			t.Errorf("expected success for provider %s", resp.Provider)
		}
	}

	if !providerNames["A"] || !providerNames["B"] {
		t.Error("expected responses from both providers A and B")
	}

	if mockRepo.GetRecordCount() != 2 {
		t.Errorf("expected 2 records in repository, got %d", mockRepo.GetRecordCount())
	}
}

func TestShippingService_BroadcastShipment_PartialFailure(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	service := NewShippingService(mockRepo)

	providerA := testutil.NewMockShippingProvider("A", "http://a.local")

	providerB := testutil.NewMockShippingProvider("B", "http://b.local")
	providerB.SetCreateShipmentFunc(func(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error) {
		return nil, errors.New("provider B failed")
	})

	service.RegisterProvider(providerA)
	service.RegisterProvider(providerB)

	request := testutil.CreateSampleShippingRequest()

	responses, err := service.BroadcastShipment(context.Background(), request)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}

	successCount := 0
	failureCount := 0
	for _, resp := range responses {
		if resp.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	if successCount != 1 {
		t.Errorf("expected 1 successful response, got %d", successCount)
	}

	if failureCount != 1 {
		t.Errorf("expected 1 failed response, got %d", failureCount)
	}

	if mockRepo.GetRecordCount() != 1 {
		t.Errorf("expected 1 record in repository (only successful), got %d", mockRepo.GetRecordCount())
	}
}

func TestShippingService_RegisterProvider(t *testing.T) {
	mockRepo := testutil.NewMockRepository()
	service := NewShippingService(mockRepo)

	if len(service.providers) != 0 {
		t.Errorf("expected 0 providers initially, got %d", len(service.providers))
	}

	provider := testutil.NewMockShippingProvider("TestProvider", "http://test.local")
	service.RegisterProvider(provider)

	if len(service.providers) != 1 {
		t.Errorf("expected 1 provider after registration, got %d", len(service.providers))
	}

	registeredProvider, exists := service.providers["TestProvider"]
	if !exists {
		t.Fatal("provider not found after registration")
	}

	if registeredProvider.GetProviderName() != "TestProvider" {
		t.Errorf("expected provider name TestProvider, got %s", registeredProvider.GetProviderName())
	}
}
