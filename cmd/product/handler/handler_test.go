package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"product/cmd/product/handler"
	"product/cmd/product/usecase"
	"product/infrastructure/log"
	"product/models"
	"product/routes"

	"github.com/gin-gonic/gin"
)

// MockProductUsecase implements usecase.ProductUsecase interface for testing.
type MockProductUsecase struct {
	GetProductByIDFn           func(ctx context.Context, productID int64) (*models.Product, error)
	GetProductCategoryByIDFn   func(ctx context.Context, productCategoryID int) (*models.ProductCategory, error)
	CreateNewProductFn         func(ctx context.Context, param *models.Product) (int64, error)
	CreateNewProductCategoryFn func(ctx context.Context, param *models.ProductCategory) (int, error)
	EditProductFn              func(ctx context.Context, param *models.Product) (*models.Product, error)
	EditProductCategoryFn      func(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error)
	DeleteProductFn            func(ctx context.Context, productID int64) error
	DeleteProductCategoryFn    func(ctx context.Context, productCategoryID int) error
	SearchProductFn            func(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error)
}

func (m *MockProductUsecase) GetProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, productID)
	}
	return &models.Product{}, nil
}

func (m *MockProductUsecase) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	if m.GetProductCategoryByIDFn != nil {
		return m.GetProductCategoryByIDFn(ctx, productCategoryID)
	}
	return &models.ProductCategory{}, nil
}

func (m *MockProductUsecase) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	if m.CreateNewProductFn != nil {
		return m.CreateNewProductFn(ctx, param)
	}
	return 1, nil
}

func (m *MockProductUsecase) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	if m.CreateNewProductCategoryFn != nil {
		return m.CreateNewProductCategoryFn(ctx, param)
	}
	return 1, nil
}

func (m *MockProductUsecase) EditProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	if m.EditProductFn != nil {
		return m.EditProductFn(ctx, param)
	}
	return param, nil
}

func (m *MockProductUsecase) EditProductCategory(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error) {
	if m.EditProductCategoryFn != nil {
		return m.EditProductCategoryFn(ctx, param)
	}
	return param, nil
}

func (m *MockProductUsecase) DeleteProduct(ctx context.Context, productID int64) error {
	if m.DeleteProductFn != nil {
		return m.DeleteProductFn(ctx, productID)
	}
	return nil
}

func (m *MockProductUsecase) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	if m.DeleteProductCategoryFn != nil {
		return m.DeleteProductCategoryFn(ctx, productCategoryID)
	}
	return nil
}

func (m *MockProductUsecase) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	if m.SearchProductFn != nil {
		return m.SearchProductFn(ctx, param)
	}
	return []models.Product{}, 0, nil
}

var _ usecase.ProductUsecase = (*MockProductUsecase)(nil)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	log.SetupLogger()
	os.Exit(m.Run())
}

func setupTestRouter(mockUC *MockProductUsecase) *gin.Engine {
	h := handler.NewProductHandler(mockUC)
	router := gin.New()
	routes.SetupRoutes(router, *h)
	return router
}

// -----------------------------------------------------------------------------
// GET /api/v1/product/:id
// -----------------------------------------------------------------------------
func TestGetProductInfo_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		GetProductByIDFn: func(ctx context.Context, productID int64) (*models.Product, error) {
			return &models.Product{
				ID:    productID,
				Name:  "Test Laptop",
				Price: 1500000,
			}, nil
		},
	}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/100", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	productObj, ok := resp["product"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'product' key in response, got %v", resp)
	}
	if fmt.Sprintf("%.0f", productObj["id"]) != "100" {
		t.Errorf("expected product ID 100, got %v", productObj["id"])
	}
}

func TestGetProductInfo_InvalidID(t *testing.T) {
	mockUC := &MockProductUsecase{}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid ID string, got %d", w.Code)
	}
}

func TestGetProductInfo_NotFound(t *testing.T) {
	mockUC := &MockProductUsecase{
		GetProductByIDFn: func(ctx context.Context, productID int64) (*models.Product, error) {
			return &models.Product{ID: 0}, nil
		},
	}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existing product, got %d", w.Code)
	}
}

func TestGetProductInfo_InternalServerError(t *testing.T) {
	mockUC := &MockProductUsecase{
		GetProductByIDFn: func(ctx context.Context, productID int64) (*models.Product, error) {
			return nil, errors.New("database connection failed")
		},
	}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 on usecase error, got %d", w.Code)
	}
}

// -----------------------------------------------------------------------------
// GET /api/v1/product_category/:id
// -----------------------------------------------------------------------------
func TestGetProductCategoryInfo_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		GetProductCategoryByIDFn: func(ctx context.Context, catID int) (*models.ProductCategory, error) {
			return &models.ProductCategory{
				ID:   catID,
				Name: "Electronics",
			}, nil
		},
	}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product_category/5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestGetProductCategoryInfo_NotFound(t *testing.T) {
	mockUC := &MockProductUsecase{
		GetProductCategoryByIDFn: func(ctx context.Context, catID int) (*models.ProductCategory, error) {
			return &models.ProductCategory{ID: 0}, nil
		},
	}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product_category/99", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// -----------------------------------------------------------------------------
// POST /api/v1/product (Management: add, edit, delete)
// -----------------------------------------------------------------------------
func TestProductManagement_Add_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		CreateNewProductFn: func(ctx context.Context, p *models.Product) (int64, error) {
			return 42, nil
		},
	}
	router := setupTestRouter(mockUC)

	param := models.ProductManagementParameter{
		Action: "add",
		Product: models.Product{
			Name:  "New Mouse",
			Price: 35000,
		},
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestProductManagement_Add_InvalidWithID(t *testing.T) {
	mockUC := &MockProductUsecase{}
	router := setupTestRouter(mockUC)

	// Action 'add' should NOT have ID set (ID must be 0)
	param := models.ProductManagementParameter{
		Action: "add",
		Product: models.Product{
			ID:   10,
			Name: "Should Fail",
		},
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 when adding product with non-zero ID, got %d", w.Code)
	}
}

func TestProductManagement_Edit_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		EditProductFn: func(ctx context.Context, p *models.Product) (*models.Product, error) {
			return p, nil
		},
	}
	router := setupTestRouter(mockUC)

	param := models.ProductManagementParameter{
		Action: "edit",
		Product: models.Product{
			ID:    5,
			Name:  "Updated Keyboard",
			Price: 120000,
		},
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestProductManagement_Edit_MissingID(t *testing.T) {
	mockUC := &MockProductUsecase{}
	router := setupTestRouter(mockUC)

	param := models.ProductManagementParameter{
		Action: "edit",
		Product: models.Product{
			ID: 0, // Missing ID
		},
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 when editing product with ID=0, got %d", w.Code)
	}
}

func TestProductManagement_Delete_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		DeleteProductFn: func(ctx context.Context, productID int64) error {
			return nil
		},
	}
	router := setupTestRouter(mockUC)

	param := models.ProductManagementParameter{
		Action: "delete",
		Product: models.Product{
			ID: 7,
		},
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestProductManagement_InvalidAction(t *testing.T) {
	mockUC := &MockProductUsecase{}
	router := setupTestRouter(mockUC)

	param := models.ProductManagementParameter{
		Action: "unknown_action",
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for unknown action, got %d", w.Code)
	}
}

// -----------------------------------------------------------------------------
// POST /api/v1/product_category (Management: add, edit, delete)
// -----------------------------------------------------------------------------
func TestProductCategoryManagement_Add_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		CreateNewProductCategoryFn: func(ctx context.Context, pc *models.ProductCategory) (int, error) {
			return 12, nil
		},
	}
	router := setupTestRouter(mockUC)

	param := models.ProductCategoryManagementParameter{
		Action: "add",
		ProductCategory: models.ProductCategory{
			Name: "Furniture",
		},
	}
	body, _ := json.Marshal(param)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/product_category", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// -----------------------------------------------------------------------------
// GET /api/v1/product/search
// -----------------------------------------------------------------------------
func TestSearchProduct_Success(t *testing.T) {
	mockUC := &MockProductUsecase{
		SearchProductFn: func(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
			items := []models.Product{
				{ID: 1, Name: "Gaming Laptop", Price: 2000000},
				{ID: 2, Name: "Office Laptop", Price: 1000000},
			}
			return items, 5, nil
		},
	}
	router := setupTestRouter(mockUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/search?name=Laptop&page=1&pageSize=2", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'data' object in search response, got %v", resp)
	}
	if fmt.Sprintf("%.0f", data["totalCount"]) != "5" {
		t.Errorf("expected totalCount 5, got %v", data["totalCount"])
	}
	if data["nextPageUrl"] == nil {
		t.Errorf("expected nextPageUrl to be present since page 1 < totalPages 3")
	}
}
