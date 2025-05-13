package usecase

import (
	"context"
	"fmt"
	"testing"

	"example-app/domain/product"
	"example-app/repository"

	"github.com/mandocaesar/mediator/pkg/mediator"
)

// mockProductDetailRepository is a mock implementation of repository.ProductDetailRepository
type mockProductDetailRepository struct {
	createDetailsFn   func(ctx context.Context, productDetail *product.ProductDetail) error
	findByProductIDFn func(ctx context.Context, productID string) (*product.ProductDetail, error)
	updateDetailsFn   func(ctx context.Context, productDetail *product.ProductDetail) error
}

func (m *mockProductDetailRepository) CreateDetails(ctx context.Context, productDetail *product.ProductDetail) error {
	if m.createDetailsFn != nil {
		return m.createDetailsFn(ctx, productDetail)
	}
	return nil
}

func (m *mockProductDetailRepository) FindByProductID(ctx context.Context, productID string) (*product.ProductDetail, error) {
	if m.findByProductIDFn != nil {
		return m.findByProductIDFn(ctx, productID)
	}
	return nil, nil
}

func (m *mockProductDetailRepository) UpdateDetails(ctx context.Context, productDetail *product.ProductDetail) error {
	if m.updateDetailsFn != nil {
		return m.updateDetailsFn(ctx, productDetail)
	}
	return nil
}

func TestProductDetailUseCase_HandleProductUpdate(t *testing.T) {
	ctx := context.Background()
	testProduct := &product.Product{
		ID:          "test_product_1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
	}

	type fields struct {
		productDetailRepo repository.ProductDetailRepository
		mediator          *mediator.Mediator
	}
	type args struct {
		ctx   context.Context
		event mediator.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful update with existing details",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{
					findByProductIDFn: func(ctx context.Context, productID string) (*product.ProductDetail, error) {
						return product.NewProductDetail(
							productID,
							"Old Manufacturer",
							"Old Category",
							0.5,
							10.0, 5.0, 2.0,
							"cm",
							map[string]string{},
						), nil
					},
					updateDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
						if productDetail.Manufacturer != "Updated Manufacturer" {
							t.Errorf("Expected Manufacturer %s, got %s", "Updated Manufacturer", productDetail.Manufacturer)
						}
						return nil
					},
				},
				mediator: mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.update",
					Payload: testProduct,
				},
			},
			wantErr: false,
		},
		{
			name: "successful update with no existing details",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{
					findByProductIDFn: func(ctx context.Context, productID string) (*product.ProductDetail, error) {
						return nil, fmt.Errorf("not found")
					},
					updateDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
						if productDetail.ProductID != testProduct.ID {
							t.Errorf("Expected ProductID %s, got %s", testProduct.ID, productDetail.ProductID)
						}
						return nil
					},
				},
				mediator: mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.update",
					Payload: testProduct,
				},
			},
			wantErr: false,
		},
		{
			name: "error updating product details",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{
					findByProductIDFn: func(ctx context.Context, productID string) (*product.ProductDetail, error) {
						return product.NewProductDetail(
							productID,
							"Old Manufacturer",
							"Old Category",
							0.5,
							10.0, 5.0, 2.0,
							"cm",
							map[string]string{},
						), nil
					},
					updateDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
						return fmt.Errorf("update failed")
					},
				},
				mediator: mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.update",
					Payload: testProduct,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid payload type",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{},
				mediator:          mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.update",
					Payload: "invalid payload",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProductDetailUseCase{
				productDetailRepo: tt.fields.productDetailRepo,
				mediator:          tt.fields.mediator,
			}
			if err := uc.HandleProductUpdate(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("ProductDetailUseCase.HandleProductUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductDetailUseCase_CreateDefaultProductDetails(t *testing.T) {
	ctx := context.Background()
	testProduct := &product.Product{
		ID:          "test_product_1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
	}

	type fields struct {
		productDetailRepo repository.ProductDetailRepository
		mediator          *mediator.Mediator
	}
	type args struct {
		ctx   context.Context
		event mediator.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful product detail creation",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{
					createDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
						if productDetail.ProductID != testProduct.ID {
							t.Errorf("Expected ProductID %s, got %s", testProduct.ID, productDetail.ProductID)
						}
						return nil
					},
				},
				mediator: mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.detail.create",
					Payload: testProduct,
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{
					createDetailsFn: func(ctx context.Context, productDetail *product.ProductDetail) error {
						return fmt.Errorf("repository error")
					},
				},
				mediator: mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.detail.create",
					Payload: testProduct,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid payload type",
			fields: fields{
				productDetailRepo: &mockProductDetailRepository{},
				mediator:          mediator.GetMediator(),
			},
			args: args{
				ctx: ctx,
				event: mediator.Event{
					Name:    "product.detail.create",
					Payload: "invalid payload",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &ProductDetailUseCase{
				productDetailRepo: tt.fields.productDetailRepo,
				mediator:          tt.fields.mediator,
			}
			if err := uc.CreateDefaultProductDetails(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("ProductDetailUseCase.CreateDefaultProductDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewProductDetailUseCase(t *testing.T) {
	// Create a mock repository
	mockRepo := &mockProductDetailRepository{}

	// Get the global mediator instance
	med := mediator.GetMediator()

	type args struct {
		productDetailRepo repository.ProductDetailRepository
	}
	tests := []struct {
		name string
		args args
		want *ProductDetailUseCase
	}{
		{
			name: "successful initialization",
			args: args{
				productDetailRepo: mockRepo,
			},
			want: &ProductDetailUseCase{
				productDetailRepo: mockRepo,
				mediator:          med,
			},
		},
		{
			name: "nil repository",
			args: args{
				productDetailRepo: nil,
			},
			want: &ProductDetailUseCase{
				productDetailRepo: nil,
				mediator:          med,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewProductDetailUseCase(tt.args.productDetailRepo)

			// Verify repository is set correctly
			if got.productDetailRepo != tt.want.productDetailRepo {
				t.Errorf("NewProductDetailUseCase().productDetailRepo = %v, want %v",
					got.productDetailRepo, tt.want.productDetailRepo)
			}

			// Verify mediator is set to global instance
			if got.mediator != tt.want.mediator {
				t.Errorf("NewProductDetailUseCase().mediator = %v, want %v",
					got.mediator, tt.want.mediator)
			}
		})
	}
}
