package product

import (
	"fmt"
	"time"
)

// Product represents the core domain model for a product
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new Product instance
func NewProduct(id, name, description string, price float64) *Product {
	now := time.Now()
	return &Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates product details
func (p *Product) Update(name, description string, price float64) {
	p.Name = name
	p.Description = description
	p.Price = price
	p.UpdatedAt = time.Now()
}

// GenerateInitialSKU creates a default SKU for the product
func (p *Product) GenerateInitialSKU() string {
	return fmt.Sprintf("SKU-%s-%s", p.ID, time.Now().Format("20060102150405"))
}

// CreateDefaultProductDetail generates a default ProductDetail for the product
func (p *Product) CreateDefaultProductDetail() *ProductDetail {
	return NewProductDetail(
		p.ID,
		"Default Manufacturer",
		"Default Category",
		0.5,
		10.0, 5.0, 2.0,
		"cm",
		map[string]string{},
	)
}
