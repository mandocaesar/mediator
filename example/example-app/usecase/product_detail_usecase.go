package usecase

import (
	"context"
	"fmt"

	"example-app/domain/product"
	"example-app/repository"

	"github.com/mandocaesar/mediator/pkg/mediator"
)

// ProductDetailUseCase handles business logic for product detail-related operations
type ProductDetailUseCase struct {
	productDetailRepo repository.ProductDetailRepository
	mediator          *mediator.Mediator
}

// NewProductDetailUseCase creates a new ProductDetailUseCase
func NewProductDetailUseCase(
	productDetailRepo repository.ProductDetailRepository,
) *ProductDetailUseCase {
	return &ProductDetailUseCase{
		productDetailRepo: productDetailRepo,
		mediator:          mediator.GetMediator(),
	}
}

// HandleProductUpdate handles product update events to update or create product details
func (uc *ProductDetailUseCase) HandleProductUpdate(ctx context.Context, event mediator.Event) error {
	product, ok := event.Payload.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid payload type for product update event")
	}

	fmt.Println("[UPDATE PRODUCT HANDLED]", event.Payload)

	// Find existing product details or create default
	existingDetails, err := uc.productDetailRepo.FindByProductID(ctx, product.ID)
	if err != nil {
		// If not found, create default details
		existingDetails = product.CreateDefaultProductDetail()
	}

	// Update product details with default or existing details
	existingDetails.Update(
		"Updated Manufacturer",
		"Updated Category",
		0.7,
		12.0, 6.0, 3.0,
		"cm",
		map[string]string{
			"Color":    "Silver",
			"Material": "Metal",
			"Updated":  "True",
		},
	)

	// Save updated product details
	err = uc.productDetailRepo.UpdateDetails(ctx, existingDetails)
	if err != nil {
		return fmt.Errorf("failed to update product details: %v", err)
	}

	// Publish an event to notify other components about product detail update
	uc.mediator.Publish(ctx, mediator.Event{
		Name:    "product.detail.updated",
		Payload: existingDetails,
	})

	return nil
}

// CreateDefaultProductDetails creates default product details for a given product
func (uc *ProductDetailUseCase) CreateDefaultProductDetails(ctx context.Context, event mediator.Event) error {
	// Extract product from event payload
	product, ok := event.Payload.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid payload type for product detail creation")
	}

	// Create default product details
	productDetail := product.CreateDefaultProductDetail()

	// Save product details
	err := uc.productDetailRepo.CreateDetails(ctx, productDetail)
	if err != nil {
		return fmt.Errorf("failed to create default product details: %v", err)
	}

	// Publish an event to notify other components about product detail creation
	uc.mediator.Publish(ctx, mediator.Event{
		Name:    "product.detail.created",
		Payload: productDetail,
	})

	return nil
}
