package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

func Init() {
	s := os.Getenv("JWT_SECRET")
	if s == "" { s = "angel-dev-secret-change-me" }
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
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { return nil, errors.New("bad method") }
		return secret, nil
	})
	if err != nil { return nil, err }
	c, ok := t.Claims.(*Claims)
	if !ok || !t.Valid { return nil, errors.New("invalid token") }
	return c, nil
}
