package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"example-app/domain/product"
	"example-app/domain/sku"

	"mediator/pkg/mediator"
)

// mockSKURepo is a mock implementation of SKURepository
type mockSKURepo struct {
	createFn          func(ctx context.Context, sku *sku.SKU) error
	updateFn          func(ctx context.Context, sku *sku.SKU) error
	findByIDFn        func(ctx context.Context, id string) (*sku.SKU, error)
	findByProductIDFn func(ctx context.Context, productID string) ([]*sku.SKU, error)
}

func (m *mockSKURepo) Create(ctx context.Context, sku *sku.SKU) error {
	if m.createFn != nil {
		return m.createFn(ctx, sku)
	}
	return nil
}

func (m *mockSKURepo) Update(ctx context.Context, sku *sku.SKU) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, sku)
	}
	return nil
}

func (m *mockSKURepo) FindByID(ctx context.Context, id string) (*sku.SKU, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockSKURepo) FindByProductID(ctx context.Context, productID string) ([]*sku.SKU, error) {
	if m.findByProductIDFn != nil {
		return m.findByProductIDFn(ctx, productID)
	}
	return nil, nil
}

func TestNewSKUUseCase(t *testing.T) {
	tests := []struct {
		name    string
		skuRepo *mockSKURepo
		wantErr bool
	}{
		{
			name:    "successful initialization",
			skuRepo: &mockSKURepo{},
			wantErr: false,
		},
		{
			name:    "nil repository",
			skuRepo: nil,
			wantErr: false, // Constructor doesn't validate nil repo
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSKUUseCase(tt.skuRepo)
			if (uc == nil) != tt.wantErr {
				t.Errorf("NewSKUUseCase() error = %v, wantErr %v", uc == nil, tt.wantErr)
			}
		})
	}
}

func TestSKUUseCase_CreateSKU(t *testing.T) {
	// Set up mock subscriber for sku.created events
	med := mediator.GetMediator()
	med.Subscribe("sku.created", func(ctx context.Context, event mediator.Event) error {
		return nil
	})

	tests := []struct {
		name        string
		skuRepo     *mockSKURepo
		productID   string
		quantity    int
		wantErr     bool
		errContains string
	}{
		{
			name: "successful creation",
			skuRepo: &mockSKURepo{
				createFn: func(ctx context.Context, sku *sku.SKU) error {
					return nil
				},
			},
			productID: "test_product_1",
			quantity:  100,
			wantErr:   false,
		},
		{
			name: "create error",
			skuRepo: &mockSKURepo{
				createFn: func(ctx context.Context, sku *sku.SKU) error {
					return errors.New("failed to create SKU")
				},
			},
			productID:   "test_product_1",
			quantity:    100,
			wantErr:     true,
			errContains: "failed to create SKU",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSKUUseCase(tt.skuRepo)
			got, err := uc.CreateSKU(context.Background(), tt.productID, tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("SKUUseCase.CreateSKU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("SKUUseCase.CreateSKU() returned nil SKU on success")
			}
			if tt.wantErr && err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("SKUUseCase.CreateSKU() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

func TestSKUUseCase_HandleSKUCreation(t *testing.T) {
	tests := []struct {
		name        string
		skuRepo     *mockSKURepo
		event       mediator.Event
		wantErr     bool
		errContains string
	}{
		{
			name: "successful handling",
			skuRepo: &mockSKURepo{
				createFn: func(ctx context.Context, sku *sku.SKU) error {
					return nil
				},
			},
			event: mediator.Event{
				Name:    "product.created",
				Payload: &product.Product{ID: "test_product_1"},
			},
			wantErr: false,
		},
		{
			name:    "invalid payload type",
			skuRepo: &mockSKURepo{},
			event: mediator.Event{
				Name:    "product.created",
				Payload: "invalid_payload",
			},
			wantErr:     true,
			errContains: "invalid payload type",
		},
		{
			name: "create error",
			skuRepo: &mockSKURepo{
				createFn: func(ctx context.Context, sku *sku.SKU) error {
					return errors.New("failed to create SKU")
				},
			},
			event: mediator.Event{
				Name:    "product.created",
				Payload: &product.Product{ID: "test_product_1"},
			},
			wantErr:     true,
			errContains: "failed to create SKU",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSKUUseCase(tt.skuRepo)
			err := uc.HandleSKUCreation(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("SKUUseCase.HandleSKUCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("SKUUseCase.HandleSKUCreation() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

func TestSKUUseCase_UpdateSKU(t *testing.T) {
	// Set up mock subscriber for sku.updated events
	med := mediator.GetMediator()
	med.Subscribe("sku.updated", func(ctx context.Context, event mediator.Event) error {
		return nil
	})

	existingSKU := sku.NewSKU("test_sku_1", "test_product_1", 100)

	tests := []struct {
		name        string
		skuRepo     *mockSKURepo
		skuID       string
		quantity    int
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update",
			skuRepo: &mockSKURepo{
				findByProductIDFn: func(ctx context.Context, productID string) ([]*sku.SKU, error) {
					return []*sku.SKU{existingSKU}, nil
				},
				updateFn: func(ctx context.Context, sku *sku.SKU) error {
					return nil
				},
			},
			skuID:    "test_sku_1",
			quantity: 200,
			wantErr:  false,
		},
		{
			name: "SKU not found",
			skuRepo: &mockSKURepo{
				findByProductIDFn: func(ctx context.Context, productID string) ([]*sku.SKU, error) {
					return nil, nil
				},
			},
			skuID:       "nonexistent_sku",
			quantity:    200,
			wantErr:     true,
			errContains: "SKU not found",
		},
		{
			name: "find error",
			skuRepo: &mockSKURepo{
				findByProductIDFn: func(ctx context.Context, productID string) ([]*sku.SKU, error) {
					return nil, errors.New("failed to find SKU")
				},
			},
			skuID:       "test_sku_1",
			quantity:    200,
			wantErr:     true,
			errContains: "SKU not found",
		},
		{
			name: "update error",
			skuRepo: &mockSKURepo{
				findByProductIDFn: func(ctx context.Context, productID string) ([]*sku.SKU, error) {
					return []*sku.SKU{existingSKU}, nil
				},
				updateFn: func(ctx context.Context, sku *sku.SKU) error {
					return errors.New("failed to update SKU")
				},
			},
			skuID:       "test_sku_1",
			quantity:    200,
			wantErr:     true,
			errContains: "failed to update SKU",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSKUUseCase(tt.skuRepo)
			err := uc.UpdateSKU(context.Background(), tt.skuID, tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("SKUUseCase.UpdateSKU() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("SKUUseCase.UpdateSKU() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}
