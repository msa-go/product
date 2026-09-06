package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"product/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	cacheKeyProductInfo         = "product:%d" // format: product:{productID} product:1
	cacheKeyProductCategoryInfo = "product_category:%d"
)

func (r *ProductRepository) GetProductByIDFromRedis(ctx context.Context, productID int64) (*models.Product, error) {
	cacheKey := fmt.Sprintf("%v:%v", cacheKeyProductInfo, productID)

	var product models.Product

	productStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.Product{}, nil
		}

		return nil, err
	}

	err = json.Unmarshal([]byte(productStr), &product)
	if err != nil {
		return nil, err
	}

	return &product, err
}

func (r *ProductRepository) GetProductCategoryByIDFromRedis(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	cacheKey := fmt.Sprintf("%v:%v", cacheKeyProductCategoryInfo, productCategoryID)

	var productCategory models.ProductCategory
	productCategoryStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.ProductCategory{}, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(productCategoryStr), &productCategory)
	if err != nil {
		return nil, err
	}

	return &productCategory, nil
}

func (r *ProductRepository) SetProductByID(ctx context.Context, product *models.Product, productID int64) error {
	cacheKey := fmt.Sprintf("%v:%v", cacheKeyProductInfo, productID)
	fmt.Println(cacheKey)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	// TODO: pick the cache TTL for a single product entry.
	// Trade-off: shorter TTL keeps price/stock fresher but increases DB
	// load on cache misses; longer TTL is cheaper but risks serving stale
	// data after a product update. Replace the placeholder below.
	expiration := 0 * time.Second

	err = r.Redis.Set(ctx, cacheKey, productJSON, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) SetProductCategoryByID(ctx, productCategory *models.ProductCategory, productCategoryID int) error {
	cacheKey := fmt.Sprintf("%v:%v", cacheKeyProductCategoryInfo, productCategoryID)

	productCategoryJSON, err := json.Marshal(productCategory)
	if err != nil {
		return err
	}

	// TODO: pick the cache TTL for a single product entry.
	expiration := 0 * time.Second

	err = r.Redis.Set(ctx, cacheKey, productCategoryJSON, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}
