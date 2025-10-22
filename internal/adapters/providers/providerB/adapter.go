package providerB

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shipping-api/internal/core/domain"
	"time"
)

type Adapter struct {
	endpoint string
	client   *http.Client
}

func NewAdapter(endpoint string) *Adapter {
	return &Adapter{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (a *Adapter) CreateShipment(ctx context.Context, request *domain.GenericShippingRequest) (*domain.ShipmentResponse, error) {
	providerReq := MapToProviderB(request)

	jsonData, err := json.Marshal(providerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rawResponse map[string]interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		rawResponse = map[string]interface{}{"raw": string(body)}
	}

	response := &domain.ShipmentResponse{
		Provider:    a.GetProviderName(),
		Success:     resp.StatusCode >= 200 && resp.StatusCode < 300,
		RawResponse: rawResponse,
	}

	if trackingID, ok := rawResponse["trackingId"].(string); ok {
		response.TrackingID = trackingID
	}
	if awb, ok := rawResponse["awb"].(string); ok {
		response.AWB = awb
	}
	if message, ok := rawResponse["message"].(string); ok {
		response.Message = message
	}

	return response, nil
}

func (a *Adapter) GetProviderName() string {
	return "B"
}

func (a *Adapter) GetEndpoint() string {
	return a.endpoint
}
