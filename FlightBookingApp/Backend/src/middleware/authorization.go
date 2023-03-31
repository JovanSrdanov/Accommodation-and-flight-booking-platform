package middleware

import (
	"FlightBookingApp/model"
	"FlightBookingApp/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.Request.Header.Get("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authentication header provided"})
		}

		tokenString := strings.Split(bearerToken, " ")[1]

		err, claims := token.VerifyToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if claims.TokenType != "access" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not an access token"})
			return
		}

		if len(ctx.Keys) == 0 {
			ctx.Keys = make(map[string]interface{})
		}

		// gets data from token, appends it to the http context
		ctx.Keys["ID"] = claims.ID
		ctx.Keys["Roles"] = claims.Roles
	}
}

// Authorization role-based authorization
// TODO Stefan: return sensible error messages
func Authorization(validRoles []model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Keys) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		roles := ctx.Keys["Roles"]
		if roles == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user has no roles"})
			return
		}

		userRoles := roles.([]model.Role)

		// checks if any of the user roles are in the valid roles group
		for _, val := range userRoles {
			if RoleMatches(val, validRoles) {
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access attempt"})
		return
	}
}

func RoleMatches(role model.Role, roles []model.Role) bool {
	returnValue := false

	for _, val := range roles {
		if role == val {
			returnValue = true
			break
		}
	}
	return returnValue
}
