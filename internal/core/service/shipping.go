package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"shipping-api/internal/core/domain"
	"shipping-api/internal/core/ports"
	"sync"
	"time"
)

type ShippingService struct {
	providers  map[string]ports.ShippingProvider
	repository ports.ShipmentRepository
}

func NewShippingService(repository ports.ShipmentRepository) *ShippingService {
	return &ShippingService{
		providers:  make(map[string]ports.ShippingProvider),
		repository: repository,
	}
}

func (s *ShippingService) RegisterProvider(provider ports.ShippingProvider) {
	s.providers[provider.GetProviderName()] = provider
}

func (s *ShippingService) ProcessShipment(ctx context.Context, request *domain.GenericShippingRequest, providerName string) (*domain.ShipmentResponse, error) {
	provider, exists := s.providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	response, err := provider.CreateShipment(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Success {
		if err := s.saveShipmentRecord(ctx, request, response); err != nil {
			return response, fmt.Errorf("failed to save shipment record: %w", err)
		}
	}

	return response, nil
}

func (s *ShippingService) BroadcastShipment(ctx context.Context, request *domain.GenericShippingRequest) ([]*domain.ShipmentResponse, error) {
	var wg sync.WaitGroup
	results := make([]*domain.ShipmentResponse, 0, len(s.providers))
	resultsChan := make(chan *domain.ShipmentResponse, len(s.providers))

	for _, provider := range s.providers {
		wg.Add(1)
		go func(p ports.ShippingProvider) {
			defer wg.Done()

			response, err := p.CreateShipment(ctx, request)
			if err != nil {
				response = &domain.ShipmentResponse{
					Provider: p.GetProviderName(),
					Success:  false,
					Message:  err.Error(),
				}
			}

			if response.Success {
				if err := s.saveShipmentRecord(ctx, request, response); err != nil {
					response.Message = fmt.Sprintf("saved failed: %v", err)
				}
			}

			resultsChan <- response
		}(provider)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		results = append(results, result)
	}

	return results, nil
}

func (s *ShippingService) saveShipmentRecord(ctx context.Context, request *domain.GenericShippingRequest, response *domain.ShipmentResponse) error {
	genericPayload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	transformedPayload, err := json.Marshal(response.RawResponse)
	if err != nil {
		return err
	}

	providerResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	record := &domain.ShipmentRecord{
		ID:                 uuid.New().String(),
		Provider:           response.Provider,
		GenericPayload:     genericPayload,
		TransformedPayload: transformedPayload,
		ProviderResponse:   providerResponse,
		Success:            response.Success,
		CreatedAt:          time.Now(),
	}

	return s.repository.Save(ctx, record)
}
