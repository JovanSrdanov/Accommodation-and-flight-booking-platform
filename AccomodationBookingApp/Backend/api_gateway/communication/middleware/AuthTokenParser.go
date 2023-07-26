package middleware

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
)

func AuthTokenParser(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// extract endpoint from context
		executedEndpoint := ctx.Request.Method + " " + ctx.Param("any")
		executedEndpoint = replaceUUIDsWithParam(executedEndpoint)
		log.Println("--> api gateway full path: ", executedEndpoint)

		allMethodsWithRoles := getProtectedMethodsWithAllowedRoles()
		allowedRoles, ok := allMethodsWithRoles[executedEndpoint]
		if !ok {
			// if a provided method is not in the accessible roles map, it means that everyone can use it
			log.Println("Endpoint: " + executedEndpoint + " not found in the list of allowed endpoints")
			ctx.Next()
			return
		}

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

		providedRole := tokenPayload.Role
		for _, role := range allowedRoles {
			if role == providedRole {
				if len(ctx.Keys) == 0 {
					ctx.Keys = make(map[string]interface{})
				}

				ctx.Keys["id"] = tokenPayload.ID
				ctx.Keys["Role"] = providedRole
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access attempt"})
		return
	}
}

func replaceUUIDsWithParam(endpoint string) string {
	// Regular expression to match UUIDs in the provided format
	uuidRegex := `[a-fA-F0-9]{24}`

	// Compile the regex pattern
	pattern := regexp.MustCompile(uuidRegex)

	// Find all matches in the input string
	matches := pattern.FindAllString(endpoint, -1)

	//log.Print("MATCHES:")
	//log.Println(matches)

	if len(matches) == 0 {
		return endpoint
	}

	// Replace each match with "ID"
	for _, _ = range matches {
		endpoint = pattern.ReplaceAllString(endpoint, "PARAM")
	}

	//log.Println("END STRING: " + endpoint)

	return endpoint
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authPathPrefix = "account-credentials/"
	const accommodationPathPrefix = "accommodation/"
	const notificationPathPrefix = "notification/"
	const ratingPathPrefix = "rating/"
	const reservationPathPrefix = "reservation/"
	const userPathPrefix = "user-profile/"

	return map[string][]model.Role{
		// authorization
		"GET /" + authPathPrefix + "username/PARAM":  {model.Guest, model.Host}, // /{username}
		"PUT /" + authPathPrefix + "change-username": {model.Guest, model.Host},
		"PUT /" + authPathPrefix + "change-password": {model.Guest, model.Host},
		"GET /" + authPathPrefix + "is-deleted":      {model.Guest, model.Host},
		// accommodation
		"GET /" + accommodationPathPrefix + "all-my":     {model.Host},
		"POST /" + accommodationPathPrefix:               {model.Host},
		"DELETE /" + accommodationPathPrefix + "by-host": {model.Host},
		// notification
		"PUT /" + notificationPathPrefix + "update-my": {model.Guest, model.Host},
		"GET /" + notificationPathPrefix + "get-my":    {model.Guest, model.Host},
		// rating
		"POST /" + ratingPathPrefix + "accommodation":         {model.Guest},
		"POST /" + ratingPathPrefix + "host":                  {model.Guest},
		"DELETE /" + ratingPathPrefix + "accommodation/PARAM": {model.Guest}, // /{accommodationId}
		"DELETE /" + ratingPathPrefix + "host/PARAM":          {model.Guest}, // /{hostId}
		// reservation
		"GET /" + "availability/all":                     {model.Host},
		"GET /" + reservationPathPrefix + "pending":      {model.Host},
		"GET /" + reservationPathPrefix + "accepted":     {model.Host},
		"GET /" + reservationPathPrefix + "accept/PARAM": {model.Host},  // /{id}
		"GET /" + reservationPathPrefix + "reject/PARAM": {model.Host},  // /{id}
		"GET /" + reservationPathPrefix + "cancel/PARAM": {model.Guest}, // /{id}
		"GET /" + reservationPathPrefix + "all/guest":    {model.Guest},
		"POST /" + reservationPathPrefix:                 {model.Guest},
		// user profile
		"PUT /" + userPathPrefix: {model.Guest, model.Host},
		"DELETE /" + "user":      {model.Guest, model.Host},
	}
}
