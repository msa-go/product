package service

import (
	"context"
	"product/cmd/product/repository"
	"product/infrastructure/log"
	"product/models"

	"github.com/sirupsen/logrus"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: productRepository,
	}
}

func (s *ProductService) DeductProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := s.ProductRepository.DeductProductStockByProductID(ctx, productID, qty)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) AddProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := s.ProductRepository.AddProductStockByProductID(ctx, productID, qty)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) GetProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	// Redis 조회
	product, err := s.ProductRepository.GetProductByIDFromRedis(ctx, productID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Errorf("s.ProductRepository.GetProductByIDFromRedis() got error %v", err)
	}

	// 존재하면 바로 반환
	if product.ID != 0 {
		return product, nil
	}

	// DB 조회
	product, err = s.ProductRepository.FindProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 존재하지 않는 상품을 캐시에 넣으면 이후 조회가 계속 빈 값을 반환하므로 건너뛴다.
	if product.ID == 0 {
		return product, nil
	}

	// // 동기 Redis 갱신
	// err = s.ProductRepository.SetProductByIDToRedis(ctx, product)
	// if err != nil {
	// 	log.Logger.WithFields(logrus.Fields{
	// 		"productID": productID,
	// 	}).Errorf("s.ProductRepository.SetProductByIDToRedis() got error %v", err)
	// }

	// 비동기 Redis 갱신
	// 요청 컨텍스트는 응답 직후 취소되므로 request_id 만 옮겨 담은 별도 컨텍스트를 사용한다.
	ctxConcurrent := context.WithValue(context.Background(), "request_id", ctx.Value("request_id"))
	go func(ctx context.Context, product *models.Product, productID int64) {
		errConcurrent := s.ProductRepository.SetProductByID(ctx, product, productID)
		if errConcurrent != nil {
			log.Logger.WithFields(logrus.Fields{
				"product": product,
			}).Errorf("s.ProductRepository.SetProductByID() got error %v", errConcurrent)
		}
	}(ctxConcurrent, product, productID)

	return product, nil
}

func (s *ProductService) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepository.FindProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (s *ProductService) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	productID, err := s.ProductRepository.InsertNewProduct(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":     param.Name,
			"category": param.CategoryID,
		}).Errorf("s.ProductRepository.InsertNewProduct() got error %v", err)
		return 0, err
	}

	return productID, nil
}

func (s *ProductService) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	productCategoryID, err := s.ProductRepository.InsertNewProductCategory(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name": param.Name,
		}).Errorf("s.ProductRepository.InsertNewProductCategory() got error %v", err)
		return 0, err
	}

	return productCategoryID, nil
}

func (s *ProductService) EditProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	product, err := s.ProductRepository.UpdateProduct(ctx, param)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) EditProductCategory(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepository.UpdateProductCategory(ctx, param)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID int64) error {
	err := s.ProductRepository.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := s.ProductRepository.DeleteProductCategory(ctx, productCategoryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := s.ProductRepository.SearchProduct(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}
