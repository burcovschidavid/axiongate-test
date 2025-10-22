package ports

import (
	"context"
	"shipping-api/internal/core/domain"
)

type ShippingProvider interface {
	CreateShipment(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error)
	GetProviderName() string
	GetEndpoint() string
}

type ShipmentRepository interface {
	Save(ctx context.Context, record *domain.ShipmentRecord) error
	FindByID(ctx context.Context, id string) (*domain.ShipmentRecord, error)
	FindByProvider(ctx context.Context, provider string, limit int) ([]*domain.ShipmentRecord, error)
}

type ShippingService interface {
	ProcessShipment(ctx context.Context, request *domain.GenericShippingRequest, providerName string) (*domain.ShipmentResponse, error)
	BroadcastShipment(ctx context.Context, request *domain.GenericShippingRequest) ([]*domain.ShipmentResponse, error)
}
