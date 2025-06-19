package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/junedAtPurplesky/purplegate-merchant-service/proto"
)

// MerchantRepository handles database operations for merchants
type MerchantRepository struct {
	pool *pgxpool.Pool
}

// NewMerchantRepository creates a new merchant repository
func NewMerchantRepository(pool *pgxpool.Pool) *MerchantRepository {
	return &MerchantRepository{pool: pool}
}

// CreateMerchant creates a new merchant in the database
func (r *MerchantRepository) CreateMerchant(ctx context.Context, name, email string) (*pb.GetMerchantResponse, error) {
	merchantID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO merchants (id, name, email, kyc_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, email, kyc_verified, created_at, updated_at
	`

	var merchant pb.GetMerchantResponse
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, query,
		merchantID,
		name,
		email,
		false, // Default kyc_verified to false
		now,
		now,
	).Scan(
		&merchant.MerchantId,
		&merchant.Name,
		&merchant.Email,
		&merchant.KycVerified,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"merchants_email_key\" (SQLSTATE 23505)" {
			return nil, fmt.Errorf("merchant with email %s already exists", email)
		}
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}

	// Convert time.Time to protobuf timestamp
	merchant.CreatedAt = timestamppb.New(createdAt)
	merchant.UpdatedAt = timestamppb.New(updatedAt)

	return &merchant, nil
}

// GetMerchant retrieves a merchant by ID
func (r *MerchantRepository) GetMerchant(ctx context.Context, merchantID string) (*pb.GetMerchantResponse, error) {
	query := `
		SELECT id, name, email, kyc_verified, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`

	var merchant pb.GetMerchantResponse
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, query, merchantID).Scan(
		&merchant.MerchantId,
		&merchant.Name,
		&merchant.Email,
		&merchant.KycVerified,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("merchant not found: %s", merchantID)
		}
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}

	// Convert time.Time to protobuf timestamp
	merchant.CreatedAt = timestamppb.New(createdAt)
	merchant.UpdatedAt = timestamppb.New(updatedAt)

	return &merchant, nil
}

// GetMerchantByEmail retrieves a merchant by email
func (r *MerchantRepository) GetMerchantByEmail(ctx context.Context, email string) (*pb.GetMerchantResponse, error) {
	query := `
		SELECT id, name, email, kyc_verified, created_at, updated_at
		FROM merchants
		WHERE email = $1
	`

	var merchant pb.GetMerchantResponse
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&merchant.MerchantId,
		&merchant.Name,
		&merchant.Email,
		&merchant.KycVerified,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("merchant not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to get merchant by email: %w", err)
	}

	// Convert time.Time to protobuf timestamp
	merchant.CreatedAt = timestamppb.New(createdAt)
	merchant.UpdatedAt = timestamppb.New(updatedAt)

	return &merchant, nil
}

// UpdateKYCStatus updates the KYC verification status of a merchant
func (r *MerchantRepository) UpdateKYCStatus(ctx context.Context, merchantID string, kycVerified bool) error {
	query := `
		UPDATE merchants
		SET kyc_verified = $2, updated_at = $3
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, merchantID, kycVerified, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update KYC status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("merchant not found: %s", merchantID)
	}

	return nil
}
