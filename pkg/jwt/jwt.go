package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("your-secret-key")

type Claims struct {
	UserID uint `json:"user-id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Snaken-TodoList",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //用签名算法创建一个新token

	signerToken, err := token.SignedString(secretKey) //使用secretKey签名
	if err != nil {
		return "", err
	}
	return signerToken, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	cllaims, ok := token.Claims.(*Claims) //获取解析后的Claims
	if !ok || !token.Valid {
		return nil, err
	}

	return cllaims, nil
}
