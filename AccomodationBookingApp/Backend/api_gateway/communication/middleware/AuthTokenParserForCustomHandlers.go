package middleware

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateToken(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authentication header provided"})
			return
		}

		accessToken := authHeader[len("Bearer "):]

		tokenPayload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if len(ctx.Keys) == 0 {
			ctx.Keys = make(map[string]interface{})
		}

		providedRole := tokenPayload.Role

		ctx.Keys["id"] = tokenPayload.ID
		ctx.Keys["Role"] = providedRole
	}
}

func Authorization(validRoles []model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Keys) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		role := ctx.Keys["Role"]
		if role == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user has no roles"})
			return
		}

		userRole := role.(model.Role)
		// checks if any of the user roles are in the valid roles group
		if roleMatches(userRole, validRoles) {
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access attempt"})
		return
	}
}

func roleMatches(role model.Role, roles []model.Role) bool {
	returnValue := false

	for _, val := range roles {
		if role == val {
			returnValue = true
			break
		}
	}
	return returnValue
}
