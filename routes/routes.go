package routes

import (
	"product/cmd/product/handler"
	"product/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productHandler handler.ProductHandler) {
	router.Use(middleware.RequestLogger())
	api := router.Group("api")

	// 상품 및 카테고리 CUD - action 필드로 처리
	api.POST("/v1/product", productHandler.ProductManagement)
}
