package main

import (
	"fmt"
	"net/http"
	"order/api/configs"
	"order/api/internal/auth"
	"order/api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listeting on port 8081")
	server.ListenAndServe()

}
