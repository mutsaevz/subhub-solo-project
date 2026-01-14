package middleware

import (
	"effective-project/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtCfg service.JWTConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "нет заголовка Authorization",
			})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "неверный заголовок в Authorization",
			})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &service.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(jwtCfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "токен недействителен либо просрочен",
			})
			return
		}

		claims, ok := token.Claims.(*service.UserClaims)

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}
		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.Role)

		ctx.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))

	for _, r := range roles {
		roleSet[r] = struct{}{}
	}

	return func(ctx *gin.Context) {
		roleVal, exists := ctx.Get("userRole")

		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "нет роли в токене",
			})
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "недопустимая роль в context",
			})
			return
		}

		if _, ok := roleSet[role]; !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "в доступе отказано",
			})
			return
		}

		ctx.Next()
	}
}
