package main

import (
	"product/cmd/product/handler"
	"product/cmd/product/repository"
	"product/cmd/product/resource"
	"product/cmd/product/service"
	"product/cmd/product/usecase"
	"product/config"
	"product/infrastructure/log"
	"product/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)
	log.SetupLogger()

	// trace.InitTracer(cfg.Observability.ServiceName, cfg.Observability.OTLPEndpoint)
	// if err != nil {
	// 	log.Logger.Fatalf("failed to init tracer: %v", err)
	// }
	// defer shutdownTracer()

	productRepository := repository.NewProductRepository(db, redis)
	productService := service.NewProductService(*productRepository)
	productUsecase := usecase.NewProductUsecase(*productService)
	productHandler := handler.NewProductHandler(*productUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *productHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
