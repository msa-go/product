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
	cacheKeyProductInfo = "product:%d" // format: product:{productID} product:1
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
