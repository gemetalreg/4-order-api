package main

import (
	"fmt"
	"net/http"
	"order/api/configs"
	"order/api/internal/auth"
	"order/api/internal/product"
	"order/api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	productRepository := product.NewProductRepository(database)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})

	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listeting on port 8081")
	server.ListenAndServe()

}
