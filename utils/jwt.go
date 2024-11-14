package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Clave secreta para firmar el token (deberías guardarla de forma segura)
var jwtKey = []byte("mi_clave_secreta")

// GenerateJWT - Genera un token JWT para un usuario
func GenerateJWT(userID int, email string) (string, error) {
	// Crear la declaración del token
	claims := &jwt.StandardClaims{
		Id:        string(userID),
		Subject:   email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expiración del token (24 horas)
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
