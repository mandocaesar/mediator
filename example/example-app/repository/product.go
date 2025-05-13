package repository

import (
	"context"
	"fmt"

	"example-app/domain/product"
)

// ProductRepository defines the interface for product persistence
type ProductRepository interface {
	// Create adds a new product to the repository
	Create(ctx context.Context, product *product.Product) error

	// FindByID retrieves a product by its unique identifier
	FindByID(ctx context.Context, id string) (*product.Product, error)

	// Update modifies an existing product in the repository
	Update(ctx context.Context, product *product.Product) error
}

// ProductDetailRepository defines the interface for product detail management
type ProductDetailRepository interface {
	// CreateDetails generates additional details for a product
	CreateDetails(ctx context.Context, productDetail *product.ProductDetail) error

	// FindByProductID retrieves product details by product ID
	FindByProductID(ctx context.Context, productID string) (*product.ProductDetail, error)

	// UpdateDetails modifies product details
	UpdateDetails(ctx context.Context, productDetail *product.ProductDetail) error
}

// In-memory implementation for demonstration
type InMemoryProductRepository struct {
	products map[string]*product.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*product.Product),
	}
}

func (r *InMemoryProductRepository) Create(ctx context.Context, p *product.Product) error {
	r.products[p.ID] = p
	return nil
}

func (r *InMemoryProductRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	p, exists := r.products[id]
	if !exists {
		return nil, fmt.Errorf("product not found")
	}
	return p, nil
}

func (r *InMemoryProductRepository) Update(ctx context.Context, p *product.Product) error {
	if _, exists := r.products[p.ID]; !exists {
		return fmt.Errorf("product not found")
	}
	r.products[p.ID] = p
	return nil
}

type InMemoryProductDetailRepository struct {
	productDetails map[string]*product.ProductDetail
}

func NewInMemoryProductDetailRepository() *InMemoryProductDetailRepository {
	return &InMemoryProductDetailRepository{
		productDetails: make(map[string]*product.ProductDetail),
	}
}

func (r *InMemoryProductDetailRepository) CreateDetails(ctx context.Context, productDetail *product.ProductDetail) error {
	r.productDetails[productDetail.ProductID] = productDetail
	return nil
}

func (r *InMemoryProductDetailRepository) FindByProductID(ctx context.Context, productID string) (*product.ProductDetail, error) {
	detail, exists := r.productDetails[productID]
	if !exists {
		return nil, fmt.Errorf("product details not found for product ID: %s", productID)
	}
	return detail, nil
}

func (r *InMemoryProductDetailRepository) UpdateDetails(ctx context.Context, productDetail *product.ProductDetail) error {
	if _, exists := r.productDetails[productDetail.ProductID]; !exists {
		return fmt.Errorf("product details not found for product ID: %s", productDetail.ProductID)
	}
	r.productDetails[productDetail.ProductID] = productDetail
	return nil
}
