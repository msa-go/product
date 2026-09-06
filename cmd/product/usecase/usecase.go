package usecase

import (
	"context"
	"product/cmd/product/service"
	"product/infrastructure/log"
	"product/models"

	"github.com/sirupsen/logrus"
)

type productUsecase struct {
	ProductService service.ProductService
}

var _ ProductUsecase = (*productUsecase)(nil)

func NewProductUsecase(productService service.ProductService) ProductUsecase {
	return &productUsecase{
		ProductService: productService,
	}
}

func (uc *productUsecase) GetProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	product, err := uc.ProductService.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *productUsecase) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.GetProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (uc *productUsecase) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	productID, err := uc.ProductService.CreateNewProduct(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":     param.Name,
			"category": param.CategoryID,
		}).Errorf("uc.ProductService.CreateNewProduct() got error %v", err)
		return 0, err
	}

	return productID, nil
}

func (uc *productUsecase) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	productCategoryID, err := uc.ProductService.CreateNewProductCategory(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name": param.Name,
		}).Errorf("uc.ProductService.CreateNewProductCategory() got error %v", err)
		return 0, err
	}

	return productCategoryID, nil
}

func (uc *productUsecase) EditProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	product, err := uc.ProductService.EditProduct(ctx, param)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *productUsecase) EditProductCategory(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.EditProductCategory(ctx, param)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, productID int64) error {
	err := uc.ProductService.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *productUsecase) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := uc.ProductService.DeleteProductCategory(ctx, productCategoryID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *productUsecase) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := uc.ProductService.SearchProduct(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}
