package auth

import (
	"fmt"
	"net/http"

	"order/api/configs"
	"order/api/pkg/req"
	"order/api/pkg/res"
)

type AuthHandler struct {
	*configs.Config
}

type AuthHandlerDeps struct {
	*configs.Config
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			return
		}
		fmt.Println(body)
		fmt.Println("Login")
		logRes := LoginResponse{
			Token: "123",
		}
		res.Json(w, logRes, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		fmt.Println(body)

		fmt.Println("Register")
		regRes := RegisterResponse{
			Token: "123",
		}
		res.Json(w, regRes, 200)

	}
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())

}
