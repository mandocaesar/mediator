package main

import (
	"context"
	"log"

	"example-app/repository"
	"example-app/usecase"

	"github.com/mandocaesar/mediator/pkg/mediator"
)

func main() {
	// Initialize mediator
	log.Println("Starting application...")
	med := mediator.GetMediator()
	log.Println("Mediator initialized successfully")

	// Create repositories
	productRepo := repository.NewInMemoryProductRepository()
	productDetailRepo := repository.NewInMemoryProductDetailRepository()
	skuRepo := repository.NewInMemorySKURepository()

	// Create use cases
	productUseCase := usecase.NewProductUseCase(
		productRepo,
		productDetailRepo,
	)

	productDetailUseCase := usecase.NewProductDetailUseCase(
		productDetailRepo,
	)

	skuUseCase := usecase.NewSKUUseCase(
		skuRepo,
	)

	// Subscribe to events
	// med.Subscribe("product.created", productUseCase.HandleProductCreation)
	// med.Subscribe("product.updated", productUseCase.HandleProductUpdate)
	// med.Subscribe("product.detail.create", productDetailUseCase.CreateDefaultProductDetails)
	med.Subscribe("product.update", productDetailUseCase.HandleProductUpdate)
	med.Subscribe("sku.created", skuUseCase.HandleSKUCreation)

	// Create a product
	ctx := context.Background()
	product, err := productUseCase.CreateProduct(
		ctx,
		"Sample Product",
		"A fantastic sample product",
		19.99,
	)
	if err != nil {
		log.Fatalf("Error creating product: %v", err)
	}
	log.Printf("Product created: %+v", product)

	// Update the product
	if err := productUseCase.UpdateProduct(ctx, product.ID,
		"Updated Product Name",
		"An even more fantastic product",
		29.99,
	); err != nil {
		log.Fatalf("Error updating product: %v", err)
	}
	log.Println("Product updated successfully")
}
