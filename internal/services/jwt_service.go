package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	cookieName = "token"
	ttl        = time.Hour
)

var signingMethod = jwt.SigningMethodHS256

type jwtService struct {
	secret []byte
}

func NewJWTService(secret string) JWTService {
	return &jwtService{secret: []byte(secret)}
}

func (s *jwtService) ParseCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	if err = cookie.Valid(); err != nil {
		return "", err
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (any, error) {
		if t.Method == nil || t.Method.Alg() != signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid: %v", token)
	}

	return claims.Subject, nil
}

func (s *jwtService) BuildCookie(userID string) (*http.Cookie, error) {
	tokenString, err := s.buildJWTString(userID)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(ttl),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie, nil
}

func (s *jwtService) buildJWTString(userID string) (string, error) {
	now := jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(signingMethod, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		NotBefore: now,
		IssuedAt:  now,
	})

	return token.SignedString(s.secret)
}
