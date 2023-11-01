package middlewares

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type OwnerClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func OwnerGenerateToken(username string, ownerSecretKey []byte) (string, error) {
	// Durasi berlakunya token
	expDate := time.Now().Add(24 * time.Hour)

	// Membuat JWT Claims
	claims := &OwnerClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expDate.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Membuat token JWT (metode HMAC)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangan token dengan secret key
	tokenStr, err := token.SignedString(ownerSecretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func OwnerVerifyToken(tokenStr string, ownerSecretKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &OwnerClaims{}, func(token *jwt.Token) (interface{}, error) {
		return ownerSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	// Validasi token
	if claims, ok := token.Claims.(*OwnerClaims); ok && token.Valid {
		return claims.Username, nil
	} else {
		return "", errors.New("invalid token")
	}
}