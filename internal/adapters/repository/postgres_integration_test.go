package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"shipping-api/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) (*PostgresRepository, func()) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/shipping_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skip("skipping integration test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("skipping integration test: database not available")
	}

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS shipment_records (
			id VARCHAR(36) PRIMARY KEY,
			provider VARCHAR(50) NOT NULL,
			generic_payload JSONB NOT NULL,
			transformed_payload JSONB NOT NULL,
			provider_response JSONB NOT NULL,
			success BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	repo := &PostgresRepository{db: db}

	cleanup := func() {
		db.Exec("DROP TABLE IF EXISTS shipment_records")
		db.Close()
	}

	return repo, cleanup
}

func TestPostgresRepository_Save(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	record := &domain.ShipmentRecord{
		ID:                 uuid.New().String(),
		Provider:           "TestProvider",
		GenericPayload:     []byte(`{"test": "data"}`),
		TransformedPayload: []byte(`{"transformed": "data"}`),
		ProviderResponse:   []byte(`{"response": "data"}`),
		Success:            true,
		CreatedAt:          time.Now(),
	}

	err := repo.Save(context.Background(), record)
	if err != nil {
		t.Fatalf("failed to save record: %v", err)
	}

	savedRecord, err := repo.FindByID(context.Background(), record.ID)
	if err != nil {
		t.Fatalf("failed to find saved record: %v", err)
	}

	if savedRecord.ID != record.ID {
		t.Errorf("expected ID %s, got %s", record.ID, savedRecord.ID)
	}

	if savedRecord.Provider != record.Provider {
		t.Errorf("expected provider %s, got %s", record.Provider, savedRecord.Provider)
	}

	if !savedRecord.Success {
		t.Error("expected success to be true")
	}

	var genericPayload map[string]interface{}
	if err := json.Unmarshal(savedRecord.GenericPayload, &genericPayload); err != nil {
		t.Fatalf("failed to unmarshal generic payload: %v", err)
	}

	if genericPayload["test"] != "data" {
		t.Error("generic payload data mismatch")
	}
}

func TestPostgresRepository_FindByID_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	record, err := repo.FindByID(context.Background(), "non-existent-id")

	if err == nil {
		t.Error("expected error for non-existent record")
	}

	if record != nil {
		t.Error("expected nil record for non-existent ID")
	}
}

func TestPostgresRepository_FindByProvider(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	providerA := "ProviderA"
	providerB := "ProviderB"

	for i := 0; i < 3; i++ {
		record := &domain.ShipmentRecord{
			ID:                 uuid.New().String(),
			Provider:           providerA,
			GenericPayload:     []byte(fmt.Sprintf(`{"test": "data%d"}`, i)),
			TransformedPayload: []byte(`{}`),
			ProviderResponse:   []byte(`{}`),
			Success:            true,
			CreatedAt:          time.Now().Add(time.Duration(-i) * time.Minute),
		}
		if err := repo.Save(context.Background(), record); err != nil {
			t.Fatalf("failed to save record: %v", err)
		}
	}

	for i := 0; i < 2; i++ {
		record := &domain.ShipmentRecord{
			ID:                 uuid.New().String(),
			Provider:           providerB,
			GenericPayload:     []byte(fmt.Sprintf(`{"test": "dataB%d"}`, i)),
			TransformedPayload: []byte(`{}`),
			ProviderResponse:   []byte(`{}`),
			Success:            true,
			CreatedAt:          time.Now(),
		}
		if err := repo.Save(context.Background(), record); err != nil {
			t.Fatalf("failed to save record: %v", err)
		}
	}

	records, err := repo.FindByProvider(context.Background(), providerA, 10)
	if err != nil {
		t.Fatalf("failed to find records by provider: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records for provider A, got %d", len(records))
	}

	for _, record := range records {
		if record.Provider != providerA {
			t.Errorf("expected provider %s, got %s", providerA, record.Provider)
		}
	}

	recordsB, err := repo.FindByProvider(context.Background(), providerB, 10)
	if err != nil {
		t.Fatalf("failed to find records by provider: %v", err)
	}

	if len(recordsB) != 2 {
		t.Errorf("expected 2 records for provider B, got %d", len(recordsB))
	}
}

func TestPostgresRepository_FindByProvider_WithLimit(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	provider := "TestProvider"

	for i := 0; i < 5; i++ {
		record := &domain.ShipmentRecord{
			ID:                 uuid.New().String(),
			Provider:           provider,
			GenericPayload:     []byte(fmt.Sprintf(`{"test": "data%d"}`, i)),
			TransformedPayload: []byte(`{}`),
			ProviderResponse:   []byte(`{}`),
			Success:            true,
			CreatedAt:          time.Now(),
		}
		if err := repo.Save(context.Background(), record); err != nil {
			t.Fatalf("failed to save record: %v", err)
		}
	}

	records, err := repo.FindByProvider(context.Background(), provider, 3)
	if err != nil {
		t.Fatalf("failed to find records: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records with limit, got %d", len(records))
	}
}

func TestPostgresRepository_SaveAndRetrieve_JSONData(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	complexPayload := map[string]interface{}{
		"weight": map[string]interface{}{
			"value": 1000,
			"unit":  "Grams",
		},
		"shipper": map[string]interface{}{
			"name":    "Test Company",
			"address": "123 Test St",
		},
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "Item 1"},
			map[string]interface{}{"id": 2, "name": "Item 2"},
		},
	}

	genericPayload, _ := json.Marshal(complexPayload)

	record := &domain.ShipmentRecord{
		ID:                 uuid.New().String(),
		Provider:           "TestProvider",
		GenericPayload:     genericPayload,
		TransformedPayload: []byte(`{"transformed": true}`),
		ProviderResponse:   []byte(`{"status": "success"}`),
		Success:            true,
		CreatedAt:          time.Now(),
	}

	if err := repo.Save(context.Background(), record); err != nil {
		t.Fatalf("failed to save record: %v", err)
	}

	savedRecord, err := repo.FindByID(context.Background(), record.ID)
	if err != nil {
		t.Fatalf("failed to find record: %v", err)
	}

	var retrievedPayload map[string]interface{}
	if err := json.Unmarshal(savedRecord.GenericPayload, &retrievedPayload); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}

	weight := retrievedPayload["weight"].(map[string]interface{})
	if weight["value"].(float64) != 1000 {
		t.Error("complex JSON data not preserved correctly")
	}

	items := retrievedPayload["items"].([]interface{})
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}
