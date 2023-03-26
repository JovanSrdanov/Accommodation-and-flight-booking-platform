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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"No authentication header provided"})
		}

		tokenString := strings.Split(bearerToken, " ")[1]

		valid, claims := token.VerifyToken(tokenString)

		if !valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Invalid authentication token"})
		}

		if len(ctx.Keys) == 0 {
			ctx.Keys = make(map[string]interface{})
		}

		// gets data from token, appends it to the http context
		ctx.Keys["ID"] = claims.ID
		ctx.Keys["Roles"] = claims.Roles
	}
}

// role-based authorization
// TODO Stefan: return sensible error messages
func Authrorization(validRoles []model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Keys) == 0 {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"request has no keys"})
		}

		rolesVal := ctx.Keys["Roles"]
		if rolesVal == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"No roles provided"})
		}

		roles := rolesVal.([]model.Role)
		validation := make(map[model.Role]int)

		for _, val := range roles {
			validation[val] = 0
		}
		// checks if the list of roles of the user matches the list of valid roles
		for _, val := range validRoles {
			if _, ok := validation[val]; !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized access attempt"})
			}
		}
	}
}