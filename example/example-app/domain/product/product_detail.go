package product

import (
	"time"
)

// ProductDetail represents additional metadata and extended information about a product
type ProductDetail struct {
	ProductID     string
	Manufacturer  string
	Category      string
	Weight        float64
	Dimensions    Dimensions
	Specifications map[string]string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Dimensions represents the physical size of a product
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Unit   string // e.g., "cm", "inches"
}

// NewProductDetail creates a new ProductDetail instance
func NewProductDetail(
	productID, manufacturer, category string, 
	weight float64, 
	length, width, height float64, 
	unit string,
	specs map[string]string,
) *ProductDetail {
	now := time.Now()
	return &ProductDetail{
		ProductID:     productID,
		Manufacturer:  manufacturer,
		Category:      category,
		Weight:        weight,
		Dimensions: Dimensions{
			Length: length,
			Width:  width,
			Height: height,
			Unit:   unit,
		},
		Specifications: specs,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Update updates the product details
func (pd *ProductDetail) Update(
	manufacturer, category string, 
	weight float64, 
	length, width, height float64, 
	unit string,
	specs map[string]string,
) {
	pd.Manufacturer = manufacturer
	pd.Category = category
	pd.Weight = weight
	pd.Dimensions = Dimensions{
		Length: length,
		Width:  width,
		Height: height,
		Unit:   unit,
	}
	pd.Specifications = specs
	pd.UpdatedAt = time.Now()
}
