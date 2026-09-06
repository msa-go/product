package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewProductRepository(database *gorm.DB, redis *redis.Client) *ProductRepository {
	return &ProductRepository{
		Database: database,
		Redis:    redis,
	}
}

// type ProductRepositoryInterface interface {
// 	FindProductByID(ctx context.Context, productID int64) (*models.Product, error)
// 	CreateNewProduct(ctx context.Context, product *models.Product) (int64, error)
// 	EditProduct(ctx context.Context, product *models.Product) (*models.Product, error)
// 	DeleteProduct(ctx context.Context, productID int64) error
// 	FindProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error)
// 	CreateNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (int, error)
// 	EditProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error)
// 	DeleteProductCategory(ctx context.Context, productCategoryID int) error
// 	SetProductByIDToRedis(ctx context.Context, product *models.Product) error
// 	SetProductCategoryByIDToRedis(ctx context.Context, productCategory *models.ProductCategory) error
// 	GetProductByIDFromRedis(ctx context.Context, productID int64) (*models.Product, error)
// 	GetProductCategoryByIDFromRedis(ctx context.Context, productCategoryID int) (*models.ProductCategory, error)
// }
