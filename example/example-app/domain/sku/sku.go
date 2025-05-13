package sku

import (
	"time"
)

// SKU represents a Stock Keeping Unit
type SKU struct {
	ID        string
	ProductID string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSKU creates a new SKU instance
func NewSKU(id, productID string, quantity int) *SKU {
	now := time.Now()
	return &SKU{
		ID:        id,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates SKU details
func (s *SKU) Update(quantity int) {
	s.Quantity = quantity
	s.UpdatedAt = time.Now()
}
