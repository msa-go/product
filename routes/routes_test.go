package routes_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"product/cmd/product/handler"
	"product/infrastructure/log"
	"product/models"
	"product/routes"

	"github.com/gin-gonic/gin"
)

type dummyUsecase struct{}

func (d *dummyUsecase) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	return &models.Product{ID: id}, nil
}
func (d *dummyUsecase) GetProductCategoryByID(ctx context.Context, id int) (*models.ProductCategory, error) {
	return &models.ProductCategory{ID: id}, nil
}
func (d *dummyUsecase) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	return 1, nil
}
func (d *dummyUsecase) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	return 1, nil
}
func (d *dummyUsecase) EditProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	return param, nil
}
func (d *dummyUsecase) EditProductCategory(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error) {
	return param, nil
}
func (d *dummyUsecase) DeleteProduct(ctx context.Context, id int64) error { return nil }
func (d *dummyUsecase) DeleteProductCategory(ctx context.Context, id int) error {
	return nil
}
func (d *dummyUsecase) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	return []models.Product{}, 0, nil
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log.SetupLogger()

	router := gin.New()
	h := handler.NewProductHandler(&dummyUsecase{})
	routes.SetupRoutes(router, *h)

	testCases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/product/1"},
		{http.MethodGet, "/api/v1/product_category/1"},
		{http.MethodGet, "/api/v1/product/search"},
		{http.MethodPost, "/api/v1/product"},
		{http.MethodPost, "/api/v1/product_category"},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			router.ServeHTTP(w, req)

			if w.Code == http.StatusNotFound {
				t.Errorf("route %s %s not registered on router", tc.method, tc.path)
			}
		})
	}
}
