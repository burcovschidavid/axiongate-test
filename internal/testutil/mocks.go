package testutil

import (
	"context"
	"shipping-api/internal/core/domain"
	"sync"
)

type MockShippingProvider struct {
	name           string
	endpoint       string
	createShipment func(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error)
}

func NewMockShippingProvider(name, endpoint string) *MockShippingProvider {
	return &MockShippingProvider{
		name:     name,
		endpoint: endpoint,
		createShipment: func(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error) {
			return &domain.ShipmentResponse{
				Provider:   name,
				Success:    true,
				TrackingID: "TRACK123",
				AWB:        "AWB123",
				Message:    "Success",
			}, nil
		},
	}
}

func (m *MockShippingProvider) CreateShipment(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error) {
	return m.createShipment(ctx, request)
}

func (m *MockShippingProvider) GetProviderName() string {
	return m.name
}

func (m *MockShippingProvider) GetEndpoint() string {
	return m.endpoint
}

func (m *MockShippingProvider) SetCreateShipmentFunc(fn func(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error)) {
	m.createShipment = fn
}

type MockRepository struct {
	records map[string]*domain.ShipmentRecord
	mu      sync.RWMutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		records: make(map[string]*domain.ShipmentRecord),
	}
}

func (m *MockRepository) Save(ctx context.Context, record *domain.ShipmentRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records[record.ID] = record
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*domain.ShipmentRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	record, exists := m.records[id]
	if !exists {
		return nil, nil
	}
	return record, nil
}

func (m *MockRepository) FindByProvider(ctx context.Context, provider string, limit int) ([]*domain.ShipmentRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var results []*domain.ShipmentRecord
	count := 0
	for _, record := range m.records {
		if record.Provider == provider && count < limit {
			results = append(results, record)
			count++
		}
	}
	return results, nil
}

func (m *MockRepository) GetRecordCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.records)
}
