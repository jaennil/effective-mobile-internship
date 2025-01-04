package service

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secretKey string
    tokenDuration time.Duration
}

type UserClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
    return &JWTManager{secretKey, tokenDuration}
}

func (m *JWTManager) GenerateJWT(user *User) (string, error) {
    claims := UserClaims {
        RegisteredClaims: jwt.RegisteredClaims {
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
        },
        Username: user.Username,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
    token, err := jwt.ParseWithClaims(
        accessToken,
        &UserClaims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected token signing method")
            }

            return []byte(m.secretKey), nil
        },
    )

    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    claims, ok := token.Claims.(*UserClaims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }

    return claims, nil
}
