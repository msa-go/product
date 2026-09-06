package usecase

import (
	"context"
	"product/cmd/product/service"
	"product/infrastructure/log"
	"product/models"

	"github.com/sirupsen/logrus"
)

type ProductUsecase struct {
	ProductService service.ProductService
}

func NewProductUsecase(productService service.ProductService) *ProductUsecase {
	return &ProductUsecase{
		ProductService: productService,
	}
}

func (uc *ProductUsecase) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
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

func (uc *ProductUsecase) EditProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	product, err := uc.ProductService.EditProduct(ctx, param)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUsecase) DeleteProduct(ctx context.Context, productID int64) error {
	err := uc.ProductService.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}
