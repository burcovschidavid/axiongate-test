package repository

import (
	"context"
	"database/sql"
	"fmt"
	"shipping-api/internal/core/domain"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Save(ctx context.Context, record *domain.ShipmentRecord) error {
	query := `
		INSERT INTO shipment_records (
			id, provider, generic_payload, transformed_payload,
			provider_response, success, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		record.ID,
		record.Provider,
		record.GenericPayload,
		record.TransformedPayload,
		record.ProviderResponse,
		record.Success,
		record.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save shipment record: %w", err)
	}

	return nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.ShipmentRecord, error) {
	query := `
		SELECT id, provider, generic_payload, transformed_payload,
			provider_response, success, created_at
		FROM shipment_records
		WHERE id = $1
	`

	record := &domain.ShipmentRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.Provider,
		&record.GenericPayload,
		&record.TransformedPayload,
		&record.ProviderResponse,
		&record.Success,
		&record.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shipment record not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find shipment record: %w", err)
	}

	return record, nil
}

func (r *PostgresRepository) FindByProvider(ctx context.Context, provider string, limit int) ([]*domain.ShipmentRecord, error) {
	query := `
		SELECT id, provider, generic_payload, transformed_payload,
			provider_response, success, created_at
		FROM shipment_records
		WHERE provider = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, provider, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query shipment records: %w", err)
	}
	defer rows.Close()

	var records []*domain.ShipmentRecord
	for rows.Next() {
		record := &domain.ShipmentRecord{}
		err := rows.Scan(
			&record.ID,
			&record.Provider,
			&record.GenericPayload,
			&record.TransformedPayload,
			&record.ProviderResponse,
			&record.Success,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shipment record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
