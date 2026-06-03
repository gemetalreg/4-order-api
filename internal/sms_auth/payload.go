package sms_auth

type SendCodeRequest struct {
	Phone string `json:"phone" validate:"required,phone"` // Добавляем валидацию телефона
}

type SendCodeResponse struct {
	SessionID string `json:"sessionId"`
}

type VerifyCodeRequest struct {
	SessionID string `json:"sessionId" validate:"required"`
	Code      string `json:"code" validate:"required,len=4,numeric"` // Код из 4 цифр
}

type VerifyCodeResponse struct {
	Token string `json:"token"`
}
