package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

var (
	ErrInvalidToken = errors.New("Invalid token")
	ErrExpiredToken = errors.New("Expired token")
)

var jwtKey = []byte("your-secret-key")

type JWTMaker interface {
	GenerateToken(userID uint) (string, error)
	VerifyToken(tokenString string) (uint, error)
}

type jwtMaker struct {
	secretKey string
}

func (maker *jwtMaker) VerifyToken(tokenString string) (uint, error) {
	return VerifyToken(tokenString)
}

func (maker *jwtMaker) GenerateToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Id:        strconv.Itoa(int(userID)),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}

func NewJWTMaker(secretKey string) JWTMaker {
	return &jwtMaker{secretKey}
}

func VerifyToken(tokenString string) (uint, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}

		// Return secret key
		return jwtKey, nil
	})

	if err != nil {
		// Check error type
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, ErrExpiredToken
			} else {
				return 0, ErrInvalidToken
			}
		}
		return 0, fmt.Errorf("Error parsing token: %v", err)
	}

	// Get claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("claims.UserID : %v", claims.UserID)
		return claims.UserID, nil
	} else {
		return 0, ErrInvalidToken
	}
}
