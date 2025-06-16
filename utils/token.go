package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    NIF      string `json:"nif"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// Gerar token JWT
func GenerateToken(userID uint, username, nif, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        NIF:      nif,
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

// GenerateRandomToken gera um token aleatório único
func GenerateRandomToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
