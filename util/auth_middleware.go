package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware intercepta y valida el token JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Falta el token de autorización"})
			return
		}
		
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parsear y validar el token usando la estructura de CustomClaims
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
			}
			return getSecretKey(), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}

		// Si el token es válido, inyectamos los claims en el contexto
		if claims, ok := token.Claims.(*CustomClaims); ok {
			c.Set("auth_type", claims.AuthType)
			c.Set("entity_id", claims.EntityID)
			
			if claims.Role != "" {
				c.Set("role", claims.Role)
			}
			if claims.Rutapod != "" {
				c.Set("rutapod", claims.Rutapod)
			}
		}

		// Continuar con la petición original
		c.Next()
	}
}
