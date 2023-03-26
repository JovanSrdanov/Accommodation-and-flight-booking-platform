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
		tokenString := strings.Split(bearerToken, " ")[1]

		valid, claims := token.VerifyToken(tokenString)

		if !valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized access")
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
func Authrorization(validRoles []model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Keys) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized access")
		}

		rolesVal := ctx.Keys["Roles"]
		if rolesVal == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized access")
		}

		roles := rolesVal.([]model.Role)
		validation := make(map[model.Role]int)

		for _, val := range roles {
			validation[val] = 0
		}
		// checks if the list of roles of the user matches the list of valid roles
		for _, val := range validRoles {
			if _, ok := validation[val]; !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized access")
			}
		}
	}
}