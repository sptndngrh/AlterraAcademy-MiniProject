package middlewares

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaimsJWT struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string, secretKey []byte) (string, error) {
	// Durasi berlakunya token
	expDate := time.Now().Add(24 * time.Hour)

	// Membuat JWT Claims
	claims := &UserClaimsJWT{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expDate.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Membuat token JWT (metode HMAC)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangan token dengan secret key
	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func VerifyTokenJWT(tokenStr string, secretKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaimsJWT{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	// Validasi token
	if claims, ok := token.Claims.(*UserClaimsJWT); ok && token.Valid {
		return claims.Username, nil
	} else {
		return "", errors.New("invalid token")
	}
}
