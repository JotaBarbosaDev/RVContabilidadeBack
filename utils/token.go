package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID   uint   `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// Gerar token JWT
func GenerateToken(userID uint, email, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Email:    email,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add((24 * time.Hour) * 10)), // Expira em 10 dias
            NotBefore: jwt.NewNumericDate(time.Now()), // Não é válido antes de agora
            IssuedAt:  jwt.NewNumericDate(time.Now()), // Emitido agora
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// Validar token JWT
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("token inválido")
}
