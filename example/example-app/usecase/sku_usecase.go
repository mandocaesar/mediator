package usecase

import (
	"context"
	"fmt"
	"time"

	"example-app/domain/product"
	"example-app/domain/sku"
	"example-app/repository"

	"github.com/mandocaesar/mediator/pkg/mediator"
)

// SKUUseCase handles business logic for SKU-related operations
type SKUUseCase struct {
	skuRepo  repository.SKURepository
	mediator *mediator.Mediator
}

// NewSKUUseCase creates a new SKUUseCase
func NewSKUUseCase(
	skuRepo repository.SKURepository,
) *SKUUseCase {
	return &SKUUseCase{
		skuRepo:  skuRepo,
		mediator: mediator.GetMediator(),
	}
}

// CreateSKU creates a new SKU for a product
func (uc *SKUUseCase) CreateSKU(ctx context.Context, productID string, quantity int) (*sku.SKU, error) {
	// Create SKU
	newSKU := sku.NewSKU(
		fmt.Sprintf("sku_%d", time.Now().UnixNano()),
		productID,
		quantity,
	)

	// Save SKU
	err := uc.skuRepo.Create(ctx, newSKU)
	if err != nil {
		return nil, err
	}

	// Publish SKU creation event
	uc.mediator.Publish(ctx, mediator.Event{
		Name:    "sku.created",
		Payload: newSKU,
	})

	return newSKU, nil
}

// HandleSKUCreation handles SKU creation events
func (uc *SKUUseCase) HandleSKUCreation(ctx context.Context, event mediator.Event) error {
	product, ok := event.Payload.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid payload type for SKU creation")
	}

	// Generate initial SKU for the product
	_, err := uc.CreateSKU(ctx, product.ID, 100) // Default initial quantity
	return err
}

// UpdateSKU updates an existing SKU
func (uc *SKUUseCase) UpdateSKU(ctx context.Context, skuID string, quantity int) error {
	// Find existing SKU
	skus, err := uc.skuRepo.FindByProductID(ctx, skuID)
	if err != nil || len(skus) == 0 {
		return fmt.Errorf("SKU not found: %s", skuID)
	}

	existingSKU := skus[0]
	existingSKU.Update(quantity)

	// Save updated SKU
	err = uc.skuRepo.Update(ctx, existingSKU)
	if err != nil {
		return err
	}

	// Publish SKU update event
	return uc.mediator.Publish(ctx, mediator.Event{
		Name:    "sku.updated",
		Payload: existingSKU,
	})
}
