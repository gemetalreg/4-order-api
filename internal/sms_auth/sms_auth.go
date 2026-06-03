package sms_auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SMSAuthService - сервис авторизации по SMS
type SMSAuthService struct {
	sessions       map[string]SessionData             // sessionId -> данные сессии
	codeExpiration time.Duration                      // время жизни кода (например, 5 минут)
	jwtSecret      string                             // секрет для JWT
	smsMock        func(phone string, code int) error // мок-функция для отправки SMS (реализация зависит от провайдера)
	rateLimiter    *RateLimiter                       // лимит запросов на отправку кодов
}

// SessionData - данные сессии, хранящиеся на сервере
type SessionData struct {
	Phone  string    // номер телефона
	Code   int       // код подтверждения
	Issued time.Time // время создания сессии
}

// NewSMSAuthService - конструктор сервиса
func NewSMSAuthService(secret string) *SMSAuthService {
	return &SMSAuthService{
		sessions:       make(map[string]SessionData),
		codeExpiration: 5 * time.Minute,
		jwtSecret:      secret,
		rateLimiter:    NewRateLimiter(3, 1*time.Minute), // max 3 attempts per minute
	}
}

// GenerateSession - генерирует сессию для указанного телефона и отправляет SMS
func (s *SMSAuthService) GenerateSession(phone string) (string, error) {
	if !s.rateLimiter.Allow(phone) {
		return "", errors.New("too many requests: try again later")
	}

	// Генерация случайного кода (4 цифры)
	code := s.generateRandomCode()

	// Генерация sessionId (UUID в виде base64)
	sessionID := s.generateSessionID()

	// Сохранение сессии в хранилище
	s.sessions[sessionID] = SessionData{
		Phone:  phone,
		Code:   code,
		Issued: time.Now(),
	}

	// Отправка SMS (мок-функция для тестов)
	if s.smsMock != nil {
		if err := s.smsMock(phone, code); err != nil {
			return "", fmt.Errorf("failed to send SMS: %v", err)
		}
	}

	return sessionID, nil
}

// VerifyCode - проверяет код подтверждения и выдает JWT
func (s *SMSAuthService) VerifyCode(sessionID string, code int) (string, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return "", errors.New("invalid session ID")
	}

	// Проверка времени жизни сессии
	if time.Since(session.Issued) > s.codeExpiration {
		return "", errors.New("code expired")
	}

	// Проверка кода
	if session.Code != code {
		return "", errors.New("invalid code")
	}

	// Удаление сессии (чтобы код не повторно использовался)
	delete(s.sessions, sessionID)

	// Генерация JWT-токена
	token, err := s.generateJWT(session.Phone)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

// generateRandomCode - генерирует случайный 4-значный код
func (s *SMSAuthService) generateRandomCode() int {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return int(b[0])%10 + (int(b[1])%10)*10 + (int(b[2])%10)*100 + (int(b[3])%10)*1000
}

// generateSessionID - генерирует случайный sessionId
func (s *SMSAuthService) generateSessionID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// generateJWT - генерирует JWT-токен для телефона
func (s *SMSAuthService) generateJWT(phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // токен действует 24 часа
	})

	return token.SignedString([]byte(s.jwtSecret))
}

// SetSMSSender - устанавливает функцию для отправки SMS (необязательно)
func (s *SMSAuthService) SetSMSSender(sender func(phone string, code int) error) {
	s.smsMock = sender
}

// RateLimiter - ограничитель запросов на отправку кодов
type RateLimiter struct {
	maxRequests int
	interval    time.Duration
	attempts    map[string]int
	lastReset   map[string]time.Time
}

func NewRateLimiter(max int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: max,
		interval:    interval,
		attempts:    make(map[string]int),
		lastReset:   make(map[string]time.Time),
	}
}

func (rl *RateLimiter) Allow(phone string) bool {
	now := time.Now()

	// Сбрасываем попытки, если прошло достаточно времени
	if lastReset, exists := rl.lastReset[phone]; exists && now.Sub(lastReset) >= rl.interval {
		rl.attempts[phone] = 0
		rl.lastReset[phone] = now
	}

	// Проверяем лимит
	if rl.attempts[phone] >= rl.maxRequests {
		return false
	}

	rl.attempts[phone]++
	return true
}
