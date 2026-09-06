package repository

import (
	"context"
	"product/models"
)

func (r *ProductRepository) InsertNewProduct(ctx context.Context, product *models.Product) (int64, error) {
	err := r.Database.WithContext(ctx).Table("product").Create(product).Error
	if err !== nil {
		return 0, err
	}

	return product.ID, nil
}

func (r *ProductRepository) InsertNewProducCategory(ctx context.Context, product *models.ProductCategory) (int64, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Create(productCategory).Error
	if err !== nil {
		return 0, err
	}

	return productCategory.ID, nil
}

// Sava() = Create or Update
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Products) (*models.Product, error) {
	err := r.Database.WithContext(ctx).Table("product").Save(product).Error
	if err != nil {
		return nil, err
	}

	return product, err
}

func (r *ProductRepository) UpdateProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Save(productCategory).Error
	if err != nil {
		return nil, err
	}

	return productCategory, err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, productID int64) error {
	err := r.Database.WithContext(ctx).Table("product").Delete(&models.Product{}, productID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := r.Database.WithContext(ctx).Table("product_category").Delete(&models.ProductCategory{}, productCategoryID).Error
	if err != nil {
		return err
	}

	return nil
}
