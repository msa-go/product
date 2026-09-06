package usecase

import (
	"context"
	"product/models"
)

type ProductUsecase interface {
	GetProductByID(ctx context.Context, productID int64) (*models.Product, error)
	GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error)
	CreateNewProduct(ctx context.Context, param *models.Product) (int64, error)
	CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error)
	EditProduct(ctx context.Context, param *models.Product) (*models.Product, error)
	EditProductCategory(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error)
	DeleteProduct(ctx context.Context, productID int64) error
	DeleteProductCategory(ctx context.Context, productCategoryID int) error
	SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error)
}
