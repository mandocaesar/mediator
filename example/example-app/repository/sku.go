package repository

import (
	"context"
	"example-app/domain/sku"
	"fmt"
	//"github.com/mandocaesar/mediator/example/example-app/domain/sku"
)

// SKURepository defines the interface for SKU persistence
type SKURepository interface {
	// Create adds a new SKU to the repository
	Create(ctx context.Context, sku *sku.SKU) error

	// FindByProductID retrieves SKUs for a specific product
	FindByProductID(ctx context.Context, productID string) ([]*sku.SKU, error)

	// Update modifies an existing SKU in the repository
	Update(ctx context.Context, sku *sku.SKU) error
}

// In-memory implementation for demonstration
type InMemorySKURepository struct {
	skus map[string]*sku.SKU
}

func NewInMemorySKURepository() *InMemorySKURepository {
	return &InMemorySKURepository{
		skus: make(map[string]*sku.SKU),
	}
}

func (r *InMemorySKURepository) Create(ctx context.Context, s *sku.SKU) error {
	r.skus[s.ID] = s
	return nil
}

func (r *InMemorySKURepository) FindByProductID(ctx context.Context, productID string) ([]*sku.SKU, error) {
	var productSKUs []*sku.SKU
	for _, s := range r.skus {
		if s.ProductID == productID {
			productSKUs = append(productSKUs, s)
		}
	}
	return productSKUs, nil
}

func (r *InMemorySKURepository) Update(ctx context.Context, s *sku.SKU) error {
	if _, exists := r.skus[s.ID]; !exists {
		return fmt.Errorf("SKU not found")
	}
	r.skus[s.ID] = s
	return nil
}
