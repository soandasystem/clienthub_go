package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Tipos de autenticación permitidos
const (
	AuthTypeUser   = "user"
	AuthTypeCourse = "course"
	AuthTypeAccess = "access_code"
)

// CustomClaims define la estructura de los datos dentro del token
type CustomClaims struct {
	AuthType string `json:"auth_type"`
	EntityID string `json:"entity_id"`
	// Campos opcionales
	Role    string `json:"role,omitempty"`
	Rutapod string `json:"rutapod,omitempty"`

	jwt.RegisteredClaims
}

// getSecretKey obtiene la clave secreta de las variables de entorno o usa una por defecto
func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("super-secreto-cambiar-en-produccion")
	}
	return []byte(secret)
}

// GenerateUserToken genera un token para el login por usuario/clave
func GenerateUserToken(userID string, role string) (string, error) {
	claims := CustomClaims{
		AuthType: AuthTypeUser,
		EntityID: userID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return createToken(claims)
}

// GenerateCourseToken genera un token para el login por apoderado/curso
func GenerateCourseToken(courseID string, rutapod string) (string, error) {
	claims := CustomClaims{
		AuthType: AuthTypeCourse,
		EntityID: courseID,
		Rutapod:  rutapod,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return createToken(claims)
}

// GenerateAccessCodeToken genera un token para el login por código de acceso
func GenerateAccessCodeToken(codeID string) (string, error) {
	claims := CustomClaims{
		AuthType: AuthTypeAccess,
		EntityID: codeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return createToken(claims)
}

// createToken es la función interna que firma el token
func createToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
