package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthTokenParser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.Next()
			return
		}

		accessToken := authHeader[len("Bearer "):]
		ctx.Set("Authorization", accessToken)

		ctx.Next()
	}
}
