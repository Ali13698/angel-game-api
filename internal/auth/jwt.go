package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

func Init() {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		// 🚨 امنیت ۱۰۰٪: در محیط پروداکشن هرگز از کلید پیش‌فرض استفاده نکنید.
		// اگر سکرت در رانفلر ست نشده باشد، سرور باید کرش کند تا جلوی جعل هویت گرفته شود.
		env := os.Getenv("APP_ENV")
		if env == "production" || env == "" {
			log.Fatal("[SECURITY FATAL] JWT_SECRET is missing in environment! Server cannot start without a secure secret.")
		}
		// فقط برای محیط توسعه (Local) از کلید موقت استفاده می‌شود
		s = "dev-only-secret-do-not-use-in-prod"
	}
	secret = []byte(s)
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func Generate(userID string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "angel",
		},
	}).SignedString(secret)
}

func Verify(token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return c, nil
}
