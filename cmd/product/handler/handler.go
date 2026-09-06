package handler

import (
	"fmt"
	"net/http"
	"product/cmd/product/usecase"
	"product/infrastructure/log"
	"product/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	ProductUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		ProductUsecase: productUsecase,
	}
}

func (h *ProductHandler) GetProductInfo(c *gin.Context) {
	productIDstr := c.Param("id")

	productID, err := strconv.ParseInt(productIDstr, 10, 64)

	// 상품 ID 파싱(str -> int64) 오류 체크
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productIDstr,
		}).Errorf("strconv.ParseInt got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Product ID",
		})

		return
	}

	product, err := h.ProductUsecase.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Errorf("h.ProductUsecase.GetProductByID() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	if product.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Info("Product ID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error_message": "Product not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) GetProductCategoryInfo(c *gin.Context) {
	productCategoryIDstr := c.Param("id")

	productCategoryID, err := strconv.Atoi(productCategoryIDstr)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productCategoryIDstr,
		}).Errorf("strconv.Atoi got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Product ID",
		})

		return
	}

	productCategory, err := h.ProductUsecase.GetProductCategoryByID(c.Request.Context(), productCategoryID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Errorf("h.ProductUsecase.GetProductCategoryByID() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	if productCategory.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Info("Product Category ID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error_message": "Product Category not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_category": productCategory,
	})
}

func (h *ProductHandler) ProductManagement(c *gin.Context) {
	var param models.ProductManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Errorf("c.ShouldBindJSON() got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})

		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing required parameter",
		})

		return
	}

	switch param.Action {
	case "add":
		// 예외처리 - 생성이기 때문에 ID는 존재하지 않아야 함 (0)
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productID, err := h.ProductUsecase.CreateNewProduct(c.Request.Context(), &param.Product)

		// 에러 처리 - 상품 생성 실패
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errof("h.ProductUsecase.CreateNewProduct() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Sucessfully create new product: %d", productID),
		})

		return

	case "edit":
		// 예외처리 - 수정이기 때문에 ID는 반드시 존재해야 함 (0보다 커야 함)
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		product, err := h.ProductUsecase.EditProduct(c.Request.Context(), &param.Product)
		// 에러 처리 - 상품 수정 실패
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errof("h.ProductUsecase.UpdateProduct() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success edit product",
			"product": product,
		})

		return

	case "delete":
		// 예외처리 - 삭제이기 때문에 ID는 반드시 존재해야 함 (0보다 커야 함)
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		err := h.ProductUsecase.DeleteProduct(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.DeleteProduct() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product %d successfully deleted!", param.ID),
		})

		return

	default:
		log.Logger.Errorf("Invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})

		return
	}
}

func (h *ProductHandler) ProductCategoryManagement(c *gin.Context) {
	var param models.ProductCategoryManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Errorf("c.ShouldBindJSON() got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})

		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing required parameter",
		})

		return
	}

	switch param.Action {
	case "add":
		// 예외처리 - 생성이기 때문에 ID는 존재하지 않아야 함 (0)
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productCategoryID, err := h.ProductUsecase.CreateNewProductCategory(c.Request.Context(), &param.ProductCategory)

		// 에러 처리 - 상품 카테고리 생성 실패
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errof("h.ProductUsecase.CreateNewProductCategory() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Sucessfully create new product category: %d", productCategoryID),
		})

		return

	case "edit":
		// 예외처리 - 수정이기 때문에 ID는 반드시 존재해야 함 (0보다 커야 함)
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product category id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productCategory, err := h.ProductUsecase.EditProductCategory(c.Request.Context(), &param.ProductCategory)
		// 에러 처리 - 상품 카테고리 수정 실패
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errof("h.ProductUsecase.EditProductCategory() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":          "Success edit product category",
			"product_category": productCategory,
		})

		return

	case "delete":
		// 예외처리 - 삭제이기 때문에 ID는 반드시 존재해야 함 (0보다 커야 함)
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product category id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		err := h.ProductUsecase.DeleteProductCategory(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.DeleteProductCategory() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product Category %d successfully deleted!", param.ID),
		})

		return

	default:
		log.Logger.Errorf("Invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})

		return
	}
}
