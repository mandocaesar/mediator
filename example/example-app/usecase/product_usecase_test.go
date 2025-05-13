package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"example-app/domain/product"

	"github.com/mandocaesar/mediator/pkg/mediator"
)

// Mock implementations
type mockProductRepo struct {
	createFn   func(ctx context.Context, product *product.Product) error
	findByIDFn func(ctx context.Context, id string) (*product.Product, error)
	updateFn   func(ctx context.Context, product *product.Product) error
}

func (m *mockProductRepo) Create(ctx context.Context, product *product.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, product)
	}
	return nil
}

func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*product.Product, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockProductRepo) Update(ctx context.Context, product *product.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, product)
	}
	return nil
}

func TestNewProductUseCase(t *testing.T) {
	productRepo := &mockProductRepo{}
	productDetailRepo := &mockProductDetailRepository{}

	uc := NewProductUseCase(productRepo, productDetailRepo)

	if uc.productRepo != productRepo {
		t.Error("NewProductUseCase() did not set product repository correctly")
	}
	if uc.productDetailRepo != productDetailRepo {
		t.Error("NewProductUseCase() did not set product detail repository correctly")
	}
	if uc.mediator != mediator.GetMediator() {
		t.Error("NewProductUseCase() did not set mediator correctly")
	}
}

func TestProductUseCase_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		productRepo    *mockProductRepo
		detailRepo     *mockProductDetailRepository
		inputName      string
		inputDesc      string
		inputPrice     float64
		wantErr        bool
		errContains    string
		checkPublished bool
	}{
		{
			name: "successful creation",
			productRepo: &mockProductRepo{
				createFn: func(ctx context.Context, product *product.Product) error {
					return nil
				},
			},
			detailRepo: &mockProductDetailRepository{
				createDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
					return nil
				},
			},
			inputName:      "Test Product",
			inputDesc:      "Test Description",
			inputPrice:     10.0,
			wantErr:        false,
			checkPublished: true,
		},
		{
			name: "product repo error",
			productRepo: &mockProductRepo{
				createFn: func(ctx context.Context, product *product.Product) error {
					return errors.New("failed to create product")
				},
			},
			detailRepo:  &mockProductDetailRepository{},
			inputName:   "Test Product",
			inputDesc:   "Test Description",
			inputPrice:  10.0,
			wantErr:     true,
			errContains: "failed to create product",
		},
		{
			name: "detail repo error",
			productRepo: &mockProductRepo{
				createFn: func(ctx context.Context, product *product.Product) error {
					return nil
				},
			},
			detailRepo: &mockProductDetailRepository{
				createDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
					return errors.New("failed to create details")
				},
			},
			inputName:   "Test Product",
			inputDesc:   "Test Description",
			inputPrice:  10.0,
			wantErr:     true,
			errContains: "failed to create details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create use case with mocks
			uc := NewProductUseCase(tt.productRepo, tt.detailRepo)

			// Execute test
			got, err := uc.CreateProduct(context.Background(), tt.inputName, tt.inputDesc, tt.inputPrice)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("CreateProduct() error = %v, want error containing %v", err, tt.errContains)
				return
			}

			// Check successful creation
			if !tt.wantErr {
				if got == nil {
					t.Error("CreateProduct() returned nil product on success")
					return
				}
				if got.Name != tt.inputName {
					t.Errorf("CreateProduct() product name = %v, want %v", got.Name, tt.inputName)
				}
				if got.Description != tt.inputDesc {
					t.Errorf("CreateProduct() product description = %v, want %v", got.Description, tt.inputDesc)
				}
				if got.Price != tt.inputPrice {
					t.Errorf("CreateProduct() product price = %v, want %v", got.Price, tt.inputPrice)
				}
			}
		})
	}
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	// Set up mock subscriber for product.updated events
	med := mediator.GetMediator()
	med.Subscribe("product.updated", func(ctx context.Context, event mediator.Event) error {
		return nil
	})

	existingProduct := product.NewProduct("test_id", "Old Name", "Old Desc", 5.0)

	tests := []struct {
		name        string
		productRepo *mockProductRepo
		detailRepo  *mockProductDetailRepository
		inputID     string
		inputName   string
		inputDesc   string
		inputPrice  float64
		wantErr     bool
		errContains string
	}{
		{
			name: "successful update",
			productRepo: &mockProductRepo{
				findByIDFn: func(ctx context.Context, id string) (*product.Product, error) {
					return existingProduct, nil
				},
				updateFn: func(ctx context.Context, product *product.Product) error {
					return nil
				},
			},
			detailRepo: &mockProductDetailRepository{
				findByProductIDFn: func(ctx context.Context, productID string) (*product.ProductDetail, error) {
					return &product.ProductDetail{}, nil
				},
			},
			inputID:    "test_id",
			inputName:  "New Name",
			inputDesc:  "New Desc",
			inputPrice: 10.0,
			wantErr:    false,
		},
		{
			name: "product not found",
			productRepo: &mockProductRepo{
				findByIDFn: func(ctx context.Context, id string) (*product.Product, error) {
					return nil, nil
				},
			},
			detailRepo:  &mockProductDetailRepository{},
			inputID:     "nonexistent_id",
			inputName:   "New Name",
			inputDesc:   "New Desc",
			inputPrice:  10.0,
			wantErr:     true,
			errContains: "product not found",
		},
		{
			name: "update error",
			productRepo: &mockProductRepo{
				findByIDFn: func(ctx context.Context, id string) (*product.Product, error) {
					return existingProduct, nil
				},
				updateFn: func(ctx context.Context, product *product.Product) error {
					return errors.New("failed to update")
				},
			},
			detailRepo:  &mockProductDetailRepository{},
			inputID:     "test_id",
			inputName:   "New Name",
			inputDesc:   "New Desc",
			inputPrice:  10.0,
			wantErr:     true,
			errContains: "failed to update",
		},
		{
			name: "detail find error",
			productRepo: &mockProductRepo{
				findByIDFn: func(ctx context.Context, id string) (*product.Product, error) {
					return existingProduct, nil
				},
				updateFn: func(ctx context.Context, product *product.Product) error {
					return nil
				},
			},
			detailRepo: &mockProductDetailRepository{
				findByProductIDFn: func(ctx context.Context, productID string) (*product.ProductDetail, error) {
					return nil, errors.New("failed to find details")
				},
			},
			inputID:     "test_id",
			inputName:   "New Name",
			inputDesc:   "New Desc",
			inputPrice:  10.0,
			wantErr:     true,
			errContains: "failed to find details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create use case with mocks
			uc := NewProductUseCase(tt.productRepo, tt.detailRepo)

			// Execute test
			err := uc.UpdateProduct(context.Background(), tt.inputID, tt.inputName, tt.inputDesc, tt.inputPrice)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("UpdateProduct() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

func TestProductUseCase_HandleProductCreation(t *testing.T) {
	tests := []struct {
		name        string
		event       mediator.Event
		wantErr     bool
		errContains string
	}{
		{
			name: "successful handling",
			event: mediator.Event{
				Name:    "product.created",
				Payload: &product.Product{ID: "test_id"},
			},
			wantErr: false,
		},
		{
			name: "invalid payload",
			event: mediator.Event{
				Name:    "product.created",
				Payload: "invalid",
			},
			wantErr:     true,
			errContains: "invalid payload type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewProductUseCase(&mockProductRepo{}, &mockProductDetailRepository{})

			err := uc.HandleProductCreation(context.Background(), tt.event)

			if (err != nil) != tt.wantErr {
				t.Errorf("HandleProductCreation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("HandleProductCreation() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

func TestProductUseCase_HandleProductUpdate(t *testing.T) {
	tests := []struct {
		name        string
		productRepo *mockProductRepo
		event       mediator.Event
		wantErr     bool
		errContains string
	}{
		{
			name: "successful handling",
			productRepo: &mockProductRepo{
				updateFn: func(ctx context.Context, product *product.Product) error {
					return nil
				},
			},
			event: mediator.Event{
				Name:    "product.updated",
				Payload: &product.Product{ID: "test_id"},
			},
			wantErr: false,
		},
		{
			name:        "invalid payload",
			productRepo: &mockProductRepo{},
			event: mediator.Event{
				Name:    "product.updated",
				Payload: "invalid",
			},
			wantErr:     true,
			errContains: "invalid payload type",
		},
		{
			name: "update error",
			productRepo: &mockProductRepo{
				updateFn: func(ctx context.Context, product *product.Product) error {
					return errors.New("update failed")
				},
			},
			event: mediator.Event{
				Name:    "product.updated",
				Payload: &product.Product{ID: "test_id"},
			},
			wantErr:     true,
			errContains: "failed to update product",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewProductUseCase(tt.productRepo, &mockProductDetailRepository{})

			err := uc.HandleProductUpdate(context.Background(), tt.event)

			if (err != nil) != tt.wantErr {
				t.Errorf("HandleProductUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("HandleProductUpdate() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}
