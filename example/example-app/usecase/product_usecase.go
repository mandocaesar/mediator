package usecase

import (
	"context"
	"fmt"
	"time"

	"example-app/domain/product"
	"example-app/repository"

	"github.com/mandocaesar/mediator/pkg/mediator"
)

// ProductUseCase handles business logic for product-related operations
type ProductUseCase struct {
	productRepo       repository.ProductRepository
	productDetailRepo repository.ProductDetailRepository
	mediator          *mediator.Mediator
}

// NewProductUseCase creates a new ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	productDetailRepo repository.ProductDetailRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:       productRepo,
		productDetailRepo: productDetailRepo,
		mediator:          mediator.GetMediator(),
	}
}

// CreateProduct creates a new product and publishes a creation event
func (uc *ProductUseCase) CreateProduct(ctx context.Context, name, description string, price float64) (*product.Product, error) {
	// Create product
	newProduct := product.NewProduct(
		fmt.Sprintf("product_%d", time.Now().UnixNano()),
		name,
		description,
		price,
	)

	// Save product
	err := uc.productRepo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	// Create product details
	productDetail := newProduct.CreateDefaultProductDetail()
	err = uc.productDetailRepo.CreateDetails(ctx, productDetail)
	if err != nil {
		return nil, err
	}

	// Publish product creation event
	uc.mediator.Publish(ctx, mediator.Event{
		Name:    "product.created",
		Payload: newProduct,
	})

	return newProduct, nil
}

// UpdateProduct updates an existing product and publishes an update event
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, productID, name, description string, price float64) error {
	// Find existing product
	existingProduct, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil || existingProduct == nil {
		return fmt.Errorf("product not found: %s", productID)
	}

	// Update product
	existingProduct.Update(name, description, price)

	// Save updated product
	err = uc.productRepo.Update(ctx, existingProduct)
	if err != nil {
		return err
	}

	// Find existing product details
	existingDetails, err := uc.productDetailRepo.FindByProductID(ctx, existingProduct.ID)
	if err != nil {
		return err
	}

	// Update product details
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

	// Publish product update event
	return uc.mediator.Publish(ctx, mediator.Event{
		Name:    "product.updated",
		Payload: existingProduct,
	})
}

// HandleProductCreation handles product creation events
func (uc *ProductUseCase) HandleProductCreation(ctx context.Context, event mediator.Event) error {
	product, ok := event.Payload.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid payload type for product creation")
	}

	// Publish event for product detail creation
	uc.mediator.Publish(ctx, mediator.Event{
		Name:    "product.detail.create",
		Payload: product,
	})

	return nil
}

// HandleProductUpdate handles product update events
func (uc *ProductUseCase) HandleProductUpdate(ctx context.Context, event mediator.Event) error {
	//log.Info("[UPDATE PRODUCT]", event)
	product, ok := event.Payload.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid payload type for product update event")
	}

	// Update product in repository
	err := uc.productRepo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}

	// Publish event for product detail update
	uc.mediator.Publish(ctx, mediator.Event{
		Name:    "product.update",
		Payload: product,
	})

	return nil
}
