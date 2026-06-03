package sms_auth

import (
	"net/http"
	"order/api/configs"
	"order/api/pkg/req"
	"order/api/pkg/res"
	"strconv"
)

// Status codes for clarity
const (
	StatusInvalidRequest = http.StatusBadRequest
	StatusServerError    = http.StatusInternalServerError
	StatusOK             = http.StatusOK
)

type AuthHandler struct {
	Service SMSAuthService
	*configs.Config
}

func NewAuthHandler(service SMSAuthService, config *configs.Config) *AuthHandler {
	return &AuthHandler{Service: service, Config: config}
}

// SendCodeHandler – POST /auth/send-code
func (h *AuthHandler) SendCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SendCodeRequest](w, r)
		if err != nil {
			res.Json(w, map[string]string{"error": "invalid request body"}, StatusInvalidRequest)
			return
		}

		if body.Phone == "" {
			res.Json(w, map[string]string{"error": "phone is required"}, StatusInvalidRequest)
			return
		}

		sessionID, err := h.Service.GenerateSession(body.Phone)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to generate session"}, StatusServerError)
			return
		}

		res.Json(w, SendCodeResponse{SessionID: sessionID}, StatusOK)
	}
}

// VerifyCodeHandler – POST /auth/verify-code
func (h *AuthHandler) VerifyCodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyCodeRequest](w, r)
		if err != nil {
			res.Json(w, map[string]string{"error": "invalid request body"}, StatusInvalidRequest)
			return
		}

		if body.Code == "" || body.SessionID == "" {
			res.Json(w, map[string]string{"error": "code and sessionID are required"}, StatusInvalidRequest)
			return
		}

		code, err := strconv.Atoi(body.Code)
		if err != nil {
			res.Json(w, map[string]string{"error": "invalid code format"}, StatusInvalidRequest)
			return
		}

		token, err := h.Service.VerifyCode(body.SessionID, code)
		if err != nil {
			res.Json(w, map[string]string{"error": "invalid or expired code"}, StatusInvalidRequest)
			return
		}

		res.Json(w, VerifyCodeResponse{Token: token}, StatusOK)
	}
}

// NewSmsAuthHandler registers routes with the provided handler (no re-creation).
func NewSmsAuthHandler(router *http.ServeMux, handler *AuthHandler) {
	router.HandleFunc("POST /auth/send-code", handler.SendCodeHandler())
	router.HandleFunc("POST /auth/verify-code", handler.VerifyCodeHandler())
}
